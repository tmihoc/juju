// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"fmt"

	"github.com/juju/errors"

	"github.com/juju/juju/core/changestream"
	"github.com/juju/juju/core/objectstore"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/eventsource"
	objectstoreerrors "github.com/juju/juju/domain/objectstore/errors"
)

// State describes retrieval and persistence methods for the objectstore.
type State interface {
	// GetMetadata returns the persistence metadata for the specified path.
	GetMetadata(ctx context.Context, path string) (objectstore.Metadata, error)
	// PutMetadata adds a new specified path for the persistence metadata.
	PutMetadata(ctx context.Context, metadata objectstore.Metadata) (objectstore.UUID, error)
	// ListMetadata returns the persistence metadata for all paths.
	ListMetadata(ctx context.Context) ([]objectstore.Metadata, error)
	// RemoveMetadata removes the specified path for the persistence metadata.
	RemoveMetadata(ctx context.Context, path string) error
	// InitialWatchStatement returns the table and the initial watch statement
	// for the persistence metadata.
	InitialWatchStatement() (string, string)
}

// WatcherFactory describes methods for creating watchers.
type WatcherFactory interface {
	// NewNamespaceWatcher returns a new namespace watcher
	// for events based on the input change mask.
	NewNamespaceWatcher(string, changestream.ChangeType, eventsource.NamespaceQuery) (watcher.StringsWatcher, error)
}

// Service provides the API for working with the objectstore.
type Service struct {
	st State
}

// NewService returns a new service reference wrapping the input state.
func NewService(st State) *Service {
	return &Service{
		st: st,
	}
}

// GetMetadata returns the persistence metadata for the specified path.
func (s *Service) GetMetadata(ctx context.Context, path string) (objectstore.Metadata, error) {
	metadata, err := s.st.GetMetadata(ctx, path)
	if err != nil {
		return objectstore.Metadata{}, errors.Annotatef(err, "retrieving metadata %s", path)
	}
	return objectstore.Metadata{
		Path:        metadata.Path,
		Hash256:     metadata.Hash256,
		Hash512_384: metadata.Hash512_384,
		Size:        metadata.Size,
	}, nil
}

// ListMetadata returns the persistence metadata for all paths.
func (s *Service) ListMetadata(ctx context.Context) ([]objectstore.Metadata, error) {
	metadata, err := s.st.ListMetadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrieving metadata: %w", err)
	}
	m := make([]objectstore.Metadata, len(metadata))
	for i, v := range metadata {
		m[i] = objectstore.Metadata{
			Path:        v.Path,
			Hash256:     v.Hash256,
			Hash512_384: v.Hash512_384,
			Size:        v.Size,
		}
	}
	return m, nil
}

// PutMetadata adds a new specified path for the persistence metadata. If any
// hash is missing, a [objectstoreerrors.ErrMissingHash] error is returned. It
// is expected that the caller supplies both hashes or none and they should be
// consistent with the object. That's the caller's responsibility.
func (s *Service) PutMetadata(ctx context.Context, metadata objectstore.Metadata) (objectstore.UUID, error) {
	// If you have one hash, you must have the other.
	if h1, h2 := metadata.Hash512_384, metadata.Hash256; h1 != "" && h2 == "" {
		return "", errors.Annotatef(objectstoreerrors.ErrMissingHash, "missing hash256")
	} else if h1 == "" && h2 != "" {
		return "", errors.Annotatef(objectstoreerrors.ErrMissingHash, "missing hash512_384")
	}

	uuid, err := s.st.PutMetadata(ctx, objectstore.Metadata{
		Hash256:     metadata.Hash256,
		Hash512_384: metadata.Hash512_384,
		Path:        metadata.Path,
		Size:        metadata.Size,
	})
	if err != nil {
		return "", errors.Annotatef(err, "adding path %s", metadata.Path)
	}

	return uuid, nil
}

// RemoveMetadata removes the specified path for the persistence metadata.
func (s *Service) RemoveMetadata(ctx context.Context, path string) error {
	err := s.st.RemoveMetadata(ctx, path)
	if err != nil {
		return errors.Annotatef(err, "removing path %s", path)
	}
	return nil
}

// WatchableService provides the API for working with the objectstore
// and the ability to create watchers.
type WatchableService struct {
	Service
	watcherFactory WatcherFactory
}

// NewWatchableService returns a new service reference wrapping the input state.
func NewWatchableService(st State, watcherFactory WatcherFactory) *WatchableService {
	return &WatchableService{
		Service: Service{
			st: st,
		},
		watcherFactory: watcherFactory,
	}
}

// Watch returns a watcher that emits the path changes that either have been
// added or removed.
func (s *WatchableService) Watch() (watcher.StringsWatcher, error) {
	table, stmt := s.st.InitialWatchStatement()
	return s.watcherFactory.NewNamespaceWatcher(
		table,
		changestream.All,
		eventsource.InitialNamespaceChanges(stmt),
	)
}
