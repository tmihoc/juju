// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package secretmigrationworker

import (
	"context"
	"runtime/debug"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/worker/v3"
	"github.com/juju/worker/v3/catacomb"
	"github.com/kr/pretty"

	coresecrets "github.com/juju/juju/core/secrets"
	"github.com/juju/juju/core/watcher"
	jujusecrets "github.com/juju/juju/secrets"
	// "github.com/juju/juju/worker/fortress"
	// "github.com/juju/juju/secrets/provider"
)

// logger is here to stop the desire of creating a package level logger.
// Don't do this, instead use the one passed as manifold config.
type logger interface{}

var _ logger = struct{}{}

// Logger represents the methods used by the worker to log information.
type Logger interface {
	Tracef(string, ...interface{})
	Debugf(string, ...interface{})
	Warningf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Criticalf(string, ...interface{})
}

// Facade instances provide a set of API for the worker to deal with secret migration task changes.
type Facade interface {
	WatchSecretBackendChanged() (watcher.NotifyWatcher, error)
	GetSecretsToMigrate() ([]coresecrets.SecretMetadataForMigration, error)
	ChangeSecretBackend(uri *coresecrets.URI, revision int, valueRef *coresecrets.ValueRef, val coresecrets.SecretData) error

	jujusecrets.JujuAPIClient
}

// Config defines the operation of the Worker.
type Config struct {
	Facade Facade
	Logger Logger
	Clock  clock.Clock

	SecretsBackendGetter func() (jujusecrets.BackendsClient, error)
}

// Validate returns an error if config cannot drive the Worker.
func (config Config) Validate() error {
	if config.Facade == nil {
		return errors.NotValidf("nil Facade")
	}
	if config.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	if config.Clock == nil {
		return errors.NotValidf("nil Clock")
	}
	if config.SecretsBackendGetter == nil {
		return errors.NotValidf("nil SecretsBackendGetter")
	}
	return nil
}

// NewWorker returns a secretmigrationworker Worker backed by config, or an error.
func NewWorker(config Config) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, errors.Trace(err)
	}

	w := &Worker{config: config}
	err := catacomb.Invoke(catacomb.Plan{
		Site: &w.catacomb,
		Work: w.loop,
	})
	return w, errors.Trace(err)
}

// Worker migrates secrets to the new backend when the model's secret backend has changed.
type Worker struct {
	catacomb catacomb.Catacomb
	config   Config

	_secretsBackends jujusecrets.BackendsClient
}

func (w *Worker) secretsBackends() (_ jujusecrets.BackendsClient, err error) {
	if w._secretsBackends == nil {
		if w._secretsBackends, err = w.config.SecretsBackendGetter(); err != nil {
			return nil, errors.Trace(err)
		}
	}
	return w._secretsBackends, nil
}

// TODO: we should do backend connection validation (ping the backend) when we change the secret backend in model config!!!
// juju model-config secret-backend=myothersecrets
// Or we should have a secret backend validator worker keeps ping the backend to make sure it's alive.
// The worker should make the backend as inactive if it's not alive.

// TODO: we should notify the secret owner (uniter or leader units) to migrate the secret to the new backend.
// Because Juju doesn't want to know the secret content at all!!!!

// TODO: user created secrets should be migrated on the controller becasue they donot have an owner unit!!

// Kill is defined on worker.Worker.
func (w *Worker) Kill() {
	w.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (w *Worker) Wait() error {
	return w.catacomb.Wait()
}

func (w *Worker) loop() (err error) {
	defer func() {
		if r := recover(); r != nil {
			w.config.Logger.Criticalf("stacktrace from panic: \n" + string(debug.Stack()))
		}
		if err != nil {
			w.config.Logger.Criticalf("secret migration worker loop failed: %v", err)
		}
	}()

	watcher, err := w.config.Facade.WatchSecretBackendChanged()
	if err != nil {
		return errors.Trace(err)
	}
	if err := w.catacomb.Add(watcher); err != nil {
		return errors.Trace(err)
	}

	for {
		select {
		case <-w.catacomb.Dying():
			return errors.Trace(w.catacomb.ErrDying())
		case _, ok := <-watcher.Changes():
			if !ok {
				return errors.New("secret backend changed watch closed")
			}

			w.config.Logger.Criticalf("secret backend has changed, now fetching secrets that needs to be migrated!!!")
			secrets, err := w.config.Facade.GetSecretsToMigrate()
			if err != nil {
				return errors.Trace(err)
			}
			if len(secrets) == 0 {
				w.config.Logger.Criticalf("no secrets needs to be migrated")
				continue
			}
			w.config.Logger.Criticalf("secrets that needs to be migrated: %+v", secrets)
			backends, err := w.secretsBackends()
			if err != nil {
				return errors.Trace(err)
			}
			w.config.Logger.Criticalf("secrets backends: %+v", backends)
			_, activeBackendID, err := backends.GetBackend(nil)
			if err != nil {
				return errors.Trace(err)
			}
			for _, md := range secrets {
				if err := w.migrateSecret(md, backends, activeBackendID); err != nil {
					return errors.Trace(err)
				}
			}
		}
	}
}

func (w *Worker) migrateSecret(
	md coresecrets.SecretMetadataForMigration,
	client jujusecrets.BackendsClient,
	activeBackendID string,
) error {
	for _, rev := range md.Revisions {
		w.config.Logger.Criticalf("migrateSecret rev(activeBackendID %q): %+v", activeBackendID, rev)
		if rev.ValueRef != nil && rev.ValueRef.BackendID == activeBackendID {
			w.config.Logger.Criticalf("THIS SHOULD NEVER HAPPEND, found secrets have already migrated!!!!")
			continue
		}

		// secretVal, err := oldBackend.GetContent(context.TODO(), rev.ValueRef.RevisionID)
		secretVal, err := client.GetRevisionContent(md.Metadata.URI, rev.Revision)
		w.config.Logger.Criticalf("migrateSecrets secretVal: %s, err %#v", pretty.Sprint(secretVal), err)
		if err != nil {
			return errors.Trace(err)
		}
		valueRef, err := client.SaveContent(md.Metadata.URI, rev.Revision, secretVal)
		w.config.Logger.Criticalf("migrateSecrets activeBackend.SaveContent(%q, %d, secretVal) err %#v", md.Metadata.URI, rev.Revision, err)
		if err != nil && !errors.Is(err, errors.NotSupported) {
			return errors.Trace(err)
		}
		var newValueRef *coresecrets.ValueRef
		data := secretVal.EncodedValues()
		if err == nil {
			// We are migrating to an external backend,
			newValueRef = &valueRef
			// The content has successfuly saved into the external backend.
			// So we won't save the content into the Juju database.
			data = nil
		}

		// var oldBackend provider.SecretsBackend
		var cleanUpExternalBackend func() error
		if rev.ValueRef != nil {
			// The old backend is an external backend.
			// Note: we have to get the old backend before we make ChangeSecretBackend facade call.
			// Because the token policy will be changed after we changed the secret's backend.
			oldBackend, _, err := client.GetBackend(&rev.ValueRef.BackendID)
			if err != nil {
				return errors.Trace(err)
			}
			cleanUpExternalBackend = func() error {
				err := oldBackend.DeleteContent(context.TODO(), rev.ValueRef.RevisionID)
				if errors.Is(err, errors.NotFound) {
					// This should never happen, but if it does, we can just ignore.
					return nil
				}
				return errors.Trace(err)
			}
		}
		if err := w.config.Facade.ChangeSecretBackend(md.Metadata.URI, rev.Revision, newValueRef, data); err != nil {
			return errors.Trace(err)
		}
		w.config.Logger.Criticalf("w.config.Facade.ChangeSecretBackend(%q, %d, %q) success", md.Metadata.URI, rev.Revision, activeBackendID)

		if cleanUpExternalBackend != nil {
			if err := cleanUpExternalBackend(); err != nil {
				return errors.Trace(err)
			}
			w.config.Logger.Criticalf("oldBackend.DeleteContent(%q) success", rev.ValueRef)
		}
	}
	return nil
}
