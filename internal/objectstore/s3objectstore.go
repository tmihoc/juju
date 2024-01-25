// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package objectstore

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"gopkg.in/tomb.v2"

	"github.com/juju/juju/core/objectstore"
)

const (
	defaultBucketName = "juju"
)

// S3ObjectStoreConfig is the configuration for the s3 object store.
type S3ObjectStoreConfig struct {
	Namespace       string
	Client          objectstore.Client
	MetadataService objectstore.ObjectStoreMetadata
	Claimer         Claimer
	Logger          Logger
	Clock           clock.Clock
}

type s3ObjectStore struct {
	baseObjectStore
	client    objectstore.Client
	namespace string
	requests  chan request
}

// NewS3ObjectStore returns a new object store worker based on the s3 backing
// storage.
func NewS3ObjectStore(ctx context.Context, cfg S3ObjectStoreConfig) (TrackedObjectStore, error) {
	s := &s3ObjectStore{
		baseObjectStore: baseObjectStore{
			claimer:         cfg.Claimer,
			metadataService: cfg.MetadataService,
			logger:          cfg.Logger,
			clock:           cfg.Clock,
		},
		client:    cfg.Client,
		namespace: cfg.Namespace,

		requests: make(chan request),
	}

	s.tomb.Go(s.loop)

	return s, nil
}

// Get returns an io.ReadCloser for data at path, namespaced to the
// model.
func (t *s3ObjectStore) Get(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	// Optimistically try to get the file from the file system. If it doesn't
	// exist, then we'll get an error, and we'll try to get it when sequencing
	// the get request with the put and remove requests.
	if reader, size, err := t.get(ctx, path); err == nil {
		return reader, size, nil
	}

	// Sequence the get request with the put and remove requests.
	response := make(chan response)
	select {
	case <-ctx.Done():
		return nil, -1, ctx.Err()
	case <-t.tomb.Dying():
		return nil, -1, tomb.ErrDying
	case t.requests <- request{
		op:       opGet,
		path:     path,
		response: response,
	}:
	}

	select {
	case <-ctx.Done():
		return nil, -1, ctx.Err()
	case <-t.tomb.Dying():
		return nil, -1, tomb.ErrDying
	case resp := <-response:
		if resp.err != nil {
			return nil, -1, fmt.Errorf("getting blob: %w", resp.err)
		}
		return resp.reader, resp.size, nil
	}
}

// Put stores data from reader at path, namespaced to the model.
func (t *s3ObjectStore) Put(ctx context.Context, path string, r io.Reader, size int64) error {
	response := make(chan response)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.tomb.Dying():
		return tomb.ErrDying
	case t.requests <- request{
		op:            opPut,
		path:          path,
		reader:        r,
		size:          size,
		hashValidator: ignoreHash,
		response:      response,
	}:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.tomb.Dying():
		return tomb.ErrDying
	case resp := <-response:
		if resp.err != nil {
			return fmt.Errorf("putting blob: %w", resp.err)
		}
		return nil
	}
}

// Put stores data from reader at path, namespaced to the model.
// It also ensures the stored data has the correct hash.
func (t *s3ObjectStore) PutAndCheckHash(ctx context.Context, path string, r io.Reader, size int64, hash string) error {
	response := make(chan response)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.tomb.Dying():
		return tomb.ErrDying
	case t.requests <- request{
		op:            opPut,
		path:          path,
		reader:        r,
		size:          size,
		hashValidator: checkHash(hash),
		response:      response,
	}:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.tomb.Dying():
		return tomb.ErrDying
	case resp := <-response:
		if resp.err != nil {
			return fmt.Errorf("putting blob and check hash: %w", resp.err)
		}
		return nil
	}
}

// Remove removes data at path, namespaced to the model.
func (t *s3ObjectStore) Remove(ctx context.Context, path string) error {
	response := make(chan response)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.tomb.Dying():
		return tomb.ErrDying
	case t.requests <- request{
		op:       opRemove,
		path:     path,
		response: response,
	}:
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.tomb.Dying():
		return tomb.ErrDying
	case resp := <-response:
		if resp.err != nil {
			return fmt.Errorf("removing blob: %w", resp.err)
		}
		return nil
	}
}

func (t *s3ObjectStore) loop() error {
	ctx, cancel := t.scopedContext()
	defer cancel()

	// Ensure that we have the base directory.
	if err := t.client.Session(ctx, func(ctx context.Context, s objectstore.Session) error {
		err := s.CreateBucket(ctx, defaultBucketName)
		if err != nil && !errors.Is(err, errors.AlreadyExists) {
			return errors.Trace(err)
		}
		return nil
	}); err != nil {
		return errors.Trace(err)
	}

	// Sequence the get request with the put, remove requests.
	for {
		select {
		case <-t.tomb.Dying():
			return tomb.ErrDying

		case req := <-t.requests:
			switch req.op {
			case opGet:
				reader, size, err := t.get(ctx, req.path)

				select {
				case <-t.tomb.Dying():
					return tomb.ErrDying

				case req.response <- response{
					reader: reader,
					size:   size,
					err:    err,
				}:
				}

			case opPut:
				select {
				case <-t.tomb.Dying():
					return tomb.ErrDying

				case req.response <- response{
					err: t.put(ctx, req.path, req.reader, req.size, req.hashValidator),
				}:
				}

			case opRemove:
				select {
				case <-t.tomb.Dying():
					return tomb.ErrDying

				case req.response <- response{
					err: t.remove(ctx, req.path),
				}:
				}
			}
		}
	}
}

func (t *s3ObjectStore) get(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	t.logger.Debugf("getting object %q from file storage", path)

	metadata, err := t.metadataService.GetMetadata(ctx, path)
	if err != nil {
		return nil, -1, fmt.Errorf("get metadata: %w", err)
	}

	var reader io.ReadCloser
	var size int64
	if err := t.client.Session(ctx, func(ctx context.Context, s objectstore.Session) error {
		var err error
		reader, size, _, err = s.GetObject(ctx, defaultBucketName, t.filePath(metadata.Hash))
		return err
	}); err != nil {
		return nil, -1, fmt.Errorf("get object: %w", err)
	}

	if metadata.Size != size {
		return nil, -1, fmt.Errorf("size mismatch for %q: expected %d, got %d", path, metadata.Size, size)
	}

	return reader, size, nil
}

func (t *s3ObjectStore) put(ctx context.Context, path string, r io.Reader, size int64, validator hashValidator) error {
	t.logger.Debugf("putting object %q to file storage", path)

	// Charms and resources are coded to use the SHA384 hash. It is possible
	// to move to the more common SHA256 hash, but that would require a
	// migration of all charms and resources during import.
	// I can only assume 384 was chosen over 256 and others, is because it's
	// not susceptible to length extension attacks? In any case, we'll
	// keep using it for now.
	fileHash := sha512.New384()

	// We need two hash sets here, because juju wants to use SHA384, but s3
	// defaults to SHA256. Luckily, we can piggyback on the writing to a tmp
	// file and create TeeReader with a MultiWriter.
	s3Hash := sha256.New()

	// We need to write this to a temp file, because if the client retries
	// then we need seek back to the beginning of the file.
	fileName, err := t.writeToTmpFile(path, io.TeeReader(r, io.MultiWriter(fileHash, s3Hash)), size)
	if err != nil {
		return errors.Trace(err)
	}

	// Encode the hashes as strings, so we can use them for file and s3 lookups.
	fileEncodedHash := hex.EncodeToString(fileHash.Sum(nil))
	s3EncodedHash := base64.StdEncoding.EncodeToString(s3Hash.Sum(nil))

	// Ensure that the hash of the file matches the expected hash.
	if expected, ok := validator(fileEncodedHash); !ok {
		return fmt.Errorf("hash mismatch for %q: expected %q, got %q: %w", path, expected, fileEncodedHash, objectstore.ErrHashMismatch)
	}

	// Lock the file with the given hash (fileEncodedHash), so that we can't
	// remove the file while we're writing it.
	return t.withLock(ctx, fileEncodedHash, func(ctx context.Context) error {
		// Persist the temporary file to the final location.
		if err := t.persistTmpFile(ctx, fileName, fileEncodedHash, s3EncodedHash, size); err != nil {
			return errors.Trace(err)
		}

		// Save the metadata for the file after we've written it. That way we
		// correctly sequence the watch events. Otherwise there is a potential
		// race where the watch event is emitted before the file is written.
		if err := t.metadataService.PutMetadata(ctx, objectstore.Metadata{
			Path: path,
			Hash: fileEncodedHash,
			Size: size,
		}); err != nil {
			return errors.Trace(err)
		}
		return nil
	})
}

func (t *s3ObjectStore) persistTmpFile(ctx context.Context, tmpFileName, fileEncodedHash, s3EncodedHash string, size int64) error {
	file, err := os.Open(tmpFileName)
	if err != nil {
		return errors.Trace(err)
	}

	if err := t.client.Session(ctx, func(ctx context.Context, s objectstore.Session) error {
		// Seek back to the beginning of the file, so that we can read it again.
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			return errors.Trace(err)
		}

		// Now that we've written the file, we can upload it to the object
		// store.
		err := s.PutObject(ctx, defaultBucketName, t.filePath(fileEncodedHash), file, s3EncodedHash)
		if err == nil {
			return nil
		}

		// If the file already exists, then we can ignore the error.
		if errors.Is(err, errors.AlreadyExists) {
			return nil
		}
		return errors.Trace(err)
	}); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (t *s3ObjectStore) remove(ctx context.Context, path string) error {
	t.logger.Debugf("removing object %q from file storage", path)

	metadata, err := t.metadataService.GetMetadata(ctx, path)
	if err != nil {
		return fmt.Errorf("get metadata: %w", err)
	}

	hash := metadata.Hash
	return t.withLock(ctx, hash, func(ctx context.Context) error {
		if err := t.metadataService.RemoveMetadata(ctx, path); err != nil {
			return fmt.Errorf("remove metadata: %w", err)
		}

		return t.client.Session(ctx, func(ctx context.Context, s objectstore.Session) error {
			err := s.DeleteObject(ctx, defaultBucketName, t.filePath(hash))
			if err == nil {
				return nil
			}
			if errors.Is(err, errors.NotFound) {
				return nil
			}
			return errors.Trace(err)
		})
	})
}

func (t *s3ObjectStore) filePath(hash string) string {
	return fmt.Sprintf("%s/%s", t.namespace, hash)
}
