// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/canonical/sqlair"
	"github.com/juju/version/v2"

	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/logger"
	coremodel "github.com/juju/juju/core/model"
	"github.com/juju/juju/core/user"
	"github.com/juju/juju/domain"
	machineerrors "github.com/juju/juju/domain/machine/errors"
	"github.com/juju/juju/domain/model"
	modelerrors "github.com/juju/juju/domain/model/errors"
	networkerrors "github.com/juju/juju/domain/network/errors"
	internaldatabase "github.com/juju/juju/internal/database"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/internal/uuid"
)

// NONEContainerType is the default container type.
var NONEContainerType = instance.NONE

// ModelState represents a type for interacting with the underlying model
// database state.
type ModelState struct {
	*domain.StateBase
	logger logger.Logger
}

// NewModelState returns a new State for interacting with the underlying model
// database state.
func NewModelState(
	factory database.TxnRunnerFactory,
	logger logger.Logger,
) *ModelState {
	return &ModelState{
		StateBase: domain.NewStateBase(factory),
		logger:    logger,
	}
}

// Create creates a new read-only model.
func (s *ModelState) Create(ctx context.Context, args model.ModelDetailArgs) error {
	db, err := s.DB()
	if err != nil {
		return errors.Capture(err)
	}

	return db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return CreateReadOnlyModel(ctx, args, s, tx)
	})
}

// Delete deletes a model.
func (s *ModelState) Delete(ctx context.Context, uuid coremodel.UUID) error {
	db, err := s.DB()
	if err != nil {
		return errors.Capture(err)
	}

	mUUID := dbUUID{UUID: uuid.String()}

	modelStmt, err := s.Prepare(`DELETE FROM model WHERE uuid = $dbUUID.uuid;`, mUUID)
	if err != nil {
		return errors.Capture(err)
	}

	// Once we get to this point, the model is hosed. We don't expect the
	// model to be in use. The model migration will reinforce the schema once
	// the migration is tried again. Failure to do that will result in the
	// model being deleted unexpected scenarios.
	modelTriggerStmt, err := s.Prepare(`DROP TRIGGER IF EXISTS trg_model_immutable_delete;`)
	if err != nil {
		return errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, modelTriggerStmt).Run()
		if errors.Is(err, sqlair.ErrNoRows) {
			return errors.New("model does not exist").Add(modelerrors.NotFound)
		} else if err != nil && !internaldatabase.IsExtendedErrorCode(err) {
			return fmt.Errorf("deleting model trigger %w", err)
		}

		var outcome sqlair.Outcome
		err = tx.Query(ctx, modelStmt, mUUID).Get(&outcome)
		if err != nil {
			return errors.Errorf("deleting readonly model information: %w", err)
		}

		if affected, err := outcome.Result().RowsAffected(); err != nil {
			return errors.Errorf("getting result from removing readonly model information: %w", err)
		} else if affected == 0 {
			return modelerrors.NotFound
		}
		return nil
	})
	if err != nil {
		return errors.Errorf("deleting model %q from model database: %w", uuid, err)
	}

	return nil
}

func (s *ModelState) getModelUUID(ctx context.Context, tx *sqlair.TX) (coremodel.UUID, error) {
	var modelUUID dbUUID
	stmt, err := s.Prepare(`SELECT &dbUUID.uuid FROM model;`, dbUUID{})
	if err != nil {
		return coremodel.UUID(""), errors.Capture(err)
	}

	err = tx.Query(ctx, stmt).Get(&modelUUID)
	if errors.Is(err, sql.ErrNoRows) {
		return coremodel.UUID(""), errors.New("model does not exist").Add(modelerrors.NotFound)
	}
	if err != nil {
		return coremodel.UUID(""), errors.Errorf("getting model uuid: %w", err)
	}

	return coremodel.UUID(modelUUID.UUID), nil
}

// GetModelConstraints returns the currently set constraints for the model.
// The following error types can be expected:
// - [modelerrors.NotFound]: when no model exists to set constraints for.
// - [modelerrors.ConstraintsNotFound]: when no model constraints have been
// set for the model.
func (s *ModelState) GetModelConstraints(ctx context.Context) (constraints.Value, error) {
	db, err := s.DB()
	if err != nil {
		return constraints.Value{}, errors.Capture(err)
	}

	selectTagStmt, err := s.Prepare(`
SELECT (ct.*) AS (&dbConstraintTag.*)
FROM constraint_tag ct
    JOIN "constraint" c ON ct.constraint_uuid = c.uuid
WHERE c.uuid = $dbConstraint.uuid`, dbConstraintTag{}, dbConstraint{})
	if err != nil {
		return constraints.Value{}, errors.Capture(err)
	}

	selectSpaceStmt, err := s.Prepare(`
SELECT (cs.*) AS (&dbConstraintSpace.*)
FROM constraint_space cs
    JOIN "constraint" c ON cs.constraint_uuid = c.uuid
WHERE c.uuid = $dbConstraint.uuid`, dbConstraintSpace{}, dbConstraint{})
	if err != nil {
		return constraints.Value{}, errors.Capture(err)
	}

	selectZoneStmt, err := s.Prepare(`
SELECT (cz.*) AS (&dbConstraintZone.*)
FROM constraint_zone cz
    JOIN "constraint" c ON cz.constraint_uuid = c.uuid
WHERE c.uuid = $dbConstraint.uuid`, dbConstraintZone{}, dbConstraint{})
	if err != nil {
		return constraints.Value{}, errors.Capture(err)
	}

	var (
		cons   dbConstraint
		tags   []dbConstraintTag
		spaces []dbConstraintSpace
		zones  []dbConstraintZone
	)
	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		cons, err = s.getModelConstraints(ctx, tx)
		if err != nil {
			return errors.Capture(err)
		}
		if cons.UUID == "" {
			// No constraint exists for the model, no furhter queries are needed.
			return nil
		}
		err = tx.Query(ctx, selectTagStmt, cons).GetAll(&tags)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("getting constraint tags: %w", err)
		}
		err = tx.Query(ctx, selectSpaceStmt, cons).GetAll(&spaces)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("getting constraint spaces: %w", err)
		}
		err = tx.Query(ctx, selectZoneStmt, cons).GetAll(&zones)
		if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
			return errors.Errorf("getting constraint zones: %w", err)
		}
		return nil
	})
	if err != nil {
		return constraints.Value{}, errors.Capture(err)
	}
	return cons.toValue(tags, spaces, zones)
}

// getConstraintUUID returns the constraint uuid that is active for the model.
// If the currently does not have any constraints then an error satisfying
// [modelerrors.ConstraintsNotFound] is returned.
func (s *ModelState) getConstrainUUID(
	ctx context.Context,
	tx *sqlair.TX,
) (string, error) {
	modelUUID, err := s.getModelUUID(ctx, tx)

	modelConstraint := dbModelConstraint{
		ModelUUID: modelUUID.String(),
	}

	stmt, err := s.Prepare(`
SELECT constraint_uuid AS &dbModelConstraint.constraint_uuid
FROM   model_constraint
WHERE  model_uuid = $dbModelConstraint.model_uuid`, modelConstraint)
	if err != nil {
		return "", errors.Capture(err)
	}

	err = tx.Query(ctx, stmt, modelConstraint).Get(&modelConstraint)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.Errorf(
			"no constraints set for model %q", modelUUID,
		).Add(modelerrors.ConstraintsNotFound)
	} else if err != nil {
		return "", errors.Errorf("getting constraint UUID for model %q: %w", modelUUID, err)
	}

	return modelConstraint.ConstraintUUID, nil
}

// getModelConstraints returns the values set in the constraints table that are
// being referenced by the model specified. If no constraints are currently set
// for the model an error satisfying [modelerrors.ConstraintsNotFound] will be
// returned.
func (s *ModelState) getModelConstraints(
	ctx context.Context,
	tx *sqlair.TX,
) (dbConstraint, error) {
	modelUUID, err := s.getModelUUID(ctx, tx)
	if err != nil {
		return dbConstraint{}, errors.Errorf("getting model uuid: %w", err)
	}

	stmt, err := s.Prepare(`
SELECT c.uuid AS &dbConstraint.uuid,
       c.arch AS &dbConstraint.arch,
       c.cpu_cores AS &dbConstraint.cpu_cores,
       c.cpu_power AS &dbConstraint.cpu_power,
       c.mem AS &dbConstraint.mem,
       c.root_disk AS &dbConstraint.root_disk,
       c.root_disk_source AS &dbConstraint.root_disk_source,
       c.instance_role AS &dbConstraint.instance_role,
       c.instance_type AS &dbConstraint.instance_type,
       ct.value AS &dbConstraint.container_type,
       c.virt_type AS &dbConstraint.virt_type,
       c.allocate_public_ip AS &dbConstraint.allocate_public_ip,
       c.image_id AS &dbConstraint.image_id
FROM   model_constraint mc
       JOIN "constraint" c ON c.uuid = mc.constraint_uuid
       LEFT JOIN container_type ct ON ct.id = c.container_type_id
WHERE  mc.model_uuid = $dbModelConstraint.model_uuid
`, dbConstraint{}, dbModelConstraint{})
	if err != nil {
		return dbConstraint{}, errors.Capture(err)
	}

	modelConstraint := dbModelConstraint{ModelUUID: modelUUID.String()}
	var constraint dbConstraint
	err = tx.Query(ctx, stmt, modelConstraint).Get(&constraint)
	if errors.Is(err, sql.ErrNoRows) {
		return dbConstraint{}, errors.Errorf(
			"no constraints set for model %q", modelUUID,
		).Add(modelerrors.ConstraintsNotFound)
	}
	if err != nil {
		return dbConstraint{}, errors.Errorf("getting model %q constraints: %w", modelUUID, err)
	}
	s.logger.Criticalf("getModelConstraints : %#v", constraint)
	return constraint, nil
}

// deleteModelConstraint deletes all constraints that are set for the provided
// model uuid. If no constraints are set for the model uuid or the model uuid
// does not exist no error is raised.
func (s *ModelState) deleteModelConstraint(
	ctx context.Context,
	tx *sqlair.TX,
) error {
	constraintUUID, err := s.getConstrainUUID(ctx, tx)
	if errors.Is(err, modelerrors.ConstraintsNotFound) {
		return nil
	} else if err != nil {
		return errors.Errorf("deleting existing model constraints: %w", err)
	}

	stmt, err := s.Prepare(`DELETE FROM model_constraint`)
	if err != nil {
		return errors.Capture(err)
	}
	err = tx.Query(ctx, stmt).Run()
	if err != nil {
		return errors.Errorf("delete constraints %q for model: %w", constraintUUID, err)
	}

	dbConstraintUUID := dbConstraintUUID{UUID: constraintUUID}

	stmt, err = s.Prepare(`
DELETE FROM constraint_tag 
WHERE constraint_uuid = $dbConstraintUUID.constraint_uuid`, dbConstraintUUID,
	)
	if err != nil {
		return errors.Capture(err)
	}

	err = tx.Query(ctx, stmt, dbConstraintUUID).Run()
	if err != nil {
		return errors.Errorf("deleting model constraint %q tags: %w", constraintUUID, err)
	}

	stmt, err = s.Prepare(`
DELETE FROM constraint_space
WHERE constraint_uuid = $dbConstraintUUID.constraint_uuid`, dbConstraintUUID,
	)
	if err != nil {
		return errors.Capture(err)
	}
	err = tx.Query(ctx, stmt, dbConstraintUUID).Run()
	if err != nil {
		return errors.Errorf("deleting model constraint %q spaces: %w", constraintUUID, err)
	}

	stmt, err = s.Prepare(`
DELETE FROM constraint_zone
WHERE constraint_uuid = $dbConstraintUUID.constraint_uuid`, dbConstraintUUID,
	)
	if err != nil {
		return errors.Capture(err)
	}
	err = tx.Query(ctx, stmt, dbConstraintUUID).Run()
	if err != nil {
		return errors.Errorf("deleting model constraint %q zones: %w", constraintUUID, err)
	}

	stmt, err = s.Prepare(`
DELETE FROM "constraint" WHERE uuid = $dbConstraintUUID.constraint_uuid`, dbConstraintUUID,
	)
	if err != nil {
		return errors.Capture(err)
	}
	err = tx.Query(ctx, stmt, dbConstraintUUID).Run()
	if err != nil {
		return errors.Errorf("deleting model constraint %q: %w", constraintUUID, err)
	}
	return nil
}

// SetModelConstraints sets the model constraints to the new values removing
// any previously set values. If the constraints container type is not set an
// error will be returned (see below). This value must be set before calling
// this method.
// The following error types can be expected:
// - [coreerrors.NotValid]: When no container type has been set in the
// constraints.
// - [networkerrors.SpaceNotFound]: when a space constraint is set but the
// space does not exist.
// - [machineerrors.InvalidContainerType]: when the container type set on the
// constraints is invalid.
// - [modelerrors.NotFound]: when no model exists to set constraints for.
func (s *ModelState) SetModelConstraints(ctx context.Context, consValue constraints.Value) error {
	db, err := s.DB()
	if err != nil {
		return errors.Capture(err)
	}

	constraintsUUID, err := uuid.NewUUID()
	if err != nil {
		return errors.Errorf("generating new model constraint uuid: %w", err)
	}

	constraintInsertValues := dbConstraintInsert{
		UUID: constraintsUUID.String(),
		Arch: sql.NullString{
			String: deref(consValue.Arch),
			Valid:  consValue.Arch != nil,
		},
		CPUCores: sql.NullInt64{
			Int64: int64(deref(consValue.CpuCores)),
			Valid: consValue.CpuCores != nil,
		},
		CPUPower: sql.NullInt64{
			Int64: int64(deref(consValue.CpuPower)),
			Valid: consValue.CpuPower != nil,
		},
		Mem: sql.NullInt64{
			Int64: int64(deref(consValue.Mem)),
			Valid: consValue.Mem != nil,
		},
		RootDisk: sql.NullInt64{
			Int64: int64(deref(consValue.RootDisk)),
			Valid: consValue.RootDisk != nil,
		},
		RootDiskSource: sql.NullString{
			String: deref(consValue.RootDiskSource),
			Valid:  consValue.RootDiskSource != nil,
		},
		InstanceRole: sql.NullString{
			String: deref(consValue.InstanceRole),
			Valid:  consValue.InstanceRole != nil,
		},
		InstanceType: sql.NullString{
			String: deref(consValue.InstanceType),
			Valid:  consValue.InstanceType != nil,
		},
		VirtType: sql.NullString{
			String: deref(consValue.VirtType),
			Valid:  consValue.VirtType != nil,
		},
		AllocatePublicIP: sql.NullBool{
			Bool:  deref(consValue.AllocatePublicIP),
			Valid: consValue.VirtType != nil,
		},
		ImageID: sql.NullString{
			String: deref(consValue.ImageID),
			Valid:  consValue.ImageID != nil,
		},
	}

	selectContainerTypeStmt, err := s.Prepare(`
SELECT &dbContainerTypeId.* FROM container_type WHERE value = $dbContainerTypeValue.value
`, dbContainerTypeId{}, dbContainerTypeValue{})
	if err != nil {
		return errors.Capture(err)
	}

	insertModelConstraintStmt, err := s.Prepare(`
INSERT INTO model_constraint (*)
VALUES ($dbModelConstraint.*)`, dbModelConstraint{})
	if err != nil {
		return errors.Capture(err)
	}

	insertConstraintStmt, err := s.Prepare(`
INSERT INTO "constraint" (*) VALUES($dbConstraintInsert.*)
`, constraintInsertValues)
	if err != nil {
		return errors.Capture(err)
	}

	return db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		modelUUID, err := s.getModelUUID(ctx, tx)
		if err != nil {
			return errors.Errorf("getting model uuid: %w", err)
		}

		err = s.deleteModelConstraint(ctx, tx)
		if err != nil {
			return errors.Errorf("deleting existing model constraints: %w", err)
		}

		if consValue.Container != nil {
			containerTypeId := dbContainerTypeId{}
			err = tx.Query(ctx, selectContainerTypeStmt, dbContainerTypeValue{
				Value: string(*consValue.Container),
			}).Get(&containerTypeId)

			if errors.Is(err, sqlair.ErrNoRows) {
				return errors.Errorf(
					"setting new constraints for model, container type %q is not valid",
					*consValue.Container,
				).Add(machineerrors.InvalidContainerType)
			} else if err != nil {
				return errors.Errorf(
					"setting new constraints for model when finding container type %q id: %w",
					string(*consValue.Container), err,
				)
			}

			constraintInsertValues.ContainerTypeId = sql.NullInt64{
				Int64: containerTypeId.Id,
				Valid: true,
			}
		}

		err = tx.Query(ctx, insertConstraintStmt, constraintInsertValues).Run()
		if err != nil {
			return errors.Errorf("setting new constraints for model: %w", err)
		}

		err = tx.Query(ctx, insertModelConstraintStmt, dbModelConstraint{
			ModelUUID:      modelUUID.String(),
			ConstraintUUID: constraintsUUID.String(),
		}).Run()
		if err != nil {
			return errors.Errorf("setting model constraints: %w", err)
		}

		if consValue.Tags != nil {
			err = s.insertConstraintTags(ctx, tx, constraintsUUID, *consValue.Tags)
			if err != nil {
				return errors.Errorf("setting constraint tags for model: %w", err)
			}
		}

		if consValue.Spaces != nil {
			err = s.insertContraintSpaces(ctx, tx, constraintsUUID, *consValue.Spaces)
			if err != nil {
				return errors.Errorf("setting constraint spaces for model: %w", err)
			}
		}

		if consValue.Zones != nil {
			err = s.insertContraintZones(ctx, tx, constraintsUUID, *consValue.Zones)
			if err != nil {
				return errors.Errorf("setting constraint zones for model: %w", err)
			}
		}
		return nil
	})
}

// insertConstraintTags is responsible for setting the specified tags for the
// supplied constraint uuid. Any previously set tags for the constraint UUID
// will not be removed. Any conflicts that exist between what has been set to be
// set will result in an error and not be handled.
func (s *ModelState) insertConstraintTags(
	ctx context.Context,
	tx *sqlair.TX,
	constraintUUID uuid.UUID,
	tags []string,
) error {
	insertConstraintTagStmt, err := s.Prepare(`
INSERT INTO constraint_tag (*)
VALUES ($dbConstraintTag.*)`, dbConstraintTag{})
	if err != nil {
		return errors.Capture(err)
	}

	if len(tags) == 0 {
		return nil
	}

	data := make([]dbConstraintTag, 0, len(tags))
	for _, tag := range tags {
		data = append(data, dbConstraintTag{
			ConstraintUUID: constraintUUID.String(),
			Tag:            tag,
		})
	}
	err = tx.Query(ctx, insertConstraintTagStmt, data).Run()
	if err != nil {
		return errors.Errorf("inserting constraint %q tags %w", constraintUUID, err)
	}
	return nil
}

// insertConstraintSpaces is responsible for setting the specified network
// spaces as constraints for the provided constraint uuid. Any previously set
// spaces for the constraint UUID will not be removed. Any conflicts that exist
// between what has been set to be set will result in an error and not be
// handled.
// If one or more of the spaces provided does not exist an error satisfying
// [networkerrors.SpaceNotFound] will be returned.
func (s *ModelState) insertContraintSpaces(
	ctx context.Context,
	tx *sqlair.TX,
	constraintUUID uuid.UUID,
	spaces []string,
) error {
	insertConstraintSpaceStmt, err := s.Prepare(`
INSERT INTO constraint_space (*)
VALUES ($dbConstraintSpace.*)`, dbConstraintSpace{})
	if err != nil {
		return errors.Capture(err)
	}

	if len(spaces) == 0 {
		return nil
	}

	data := make([]dbConstraintSpace, 0, len(spaces))
	for _, space := range spaces {
		data = append(data, dbConstraintSpace{
			ConstraintUUID: constraintUUID.String(),
			Space:          space,
		})
	}
	err = tx.Query(ctx, insertConstraintSpaceStmt, data).Run()
	if internaldatabase.IsErrConstraintForeignKey(err) {
		return errors.Errorf(
			"inserting constraints %q spaces, space(s) %v does not exist",
			constraintUUID,
		).Add(networkerrors.SpaceNotFound)
	}
	if err != nil {
		return errors.Errorf("inserting constraint %q space(s): %w", err)
	}
	return nil
}

// insertConstraintZones is responsible for setting the specified zones as
// constraints on the provided constraint uuid. Any previously set zones for the
// constraint UUID will not be removed. Any conflicts that exist between what
// has been set to be set will result in an error and not be handled.
func (s *ModelState) insertContraintZones(
	ctx context.Context,
	tx *sqlair.TX,
	constraintUUID uuid.UUID,
	zones []string,
) error {
	insertConstraintZoneStmt, err := s.Prepare(`
INSERT INTO constraint_zone (*)
VALUES ($dbConstraintZone.*)`, dbConstraintZone{})
	if err != nil {
		return errors.Capture(err)
	}

	if len(zones) == 0 {
		return nil
	}

	data := make([]dbConstraintZone, 0, len(zones))
	for _, zone := range zones {
		data = append(data, dbConstraintZone{
			ConstraintUUID: constraintUUID.String(),
			Zone:           zone,
		})
	}
	err = tx.Query(ctx, insertConstraintZoneStmt, data).Run()
	if err != nil {
		return errors.Errorf("inserting constraint zone: %w", err)
	}
	return nil
}

// GetModel returns a read-only model information that has been set in the
// database. If no model has been set then an error satisfying
// [modelerrors.NotFound] is returned.
func (s *ModelState) GetModel(ctx context.Context) (coremodel.ModelInfo, error) {
	db, err := s.DB()
	if err != nil {
		return coremodel.ModelInfo{}, errors.Capture(err)
	}

	m := dbReadOnlyModel{}
	stmt, err := s.Prepare(`SELECT &dbReadOnlyModel.* FROM model`, m)
	if err != nil {
		return coremodel.ModelInfo{}, errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt).Get(&m)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("model does not exist").Add(modelerrors.NotFound)
		}
		return err
	})

	if err != nil {
		return coremodel.ModelInfo{}, errors.Errorf(
			"getting model read only information: %w", err,
		)
	}

	model := coremodel.ModelInfo{
		UUID:              coremodel.UUID(m.UUID),
		Name:              m.Name,
		Type:              coremodel.ModelType(m.Type),
		Cloud:             m.Cloud,
		CloudType:         m.CloudType,
		CloudRegion:       m.CloudRegion,
		CredentialName:    m.CredentialName,
		IsControllerModel: m.IsControllerModel,
	}

	if owner := m.CredentialOwner; owner != "" {
		model.CredentialOwner, err = user.NewName(owner)
		if err != nil {
			return coremodel.ModelInfo{}, errors.Errorf(
				"parsing model %q owner username %q: %w",
				m.UUID, owner, err,
			)
		}
	} else {
		s.logger.Infof("model %s: cloud credential owner name is empty", model.Name)
	}

	var agentVersion string
	if m.TargetAgentVersion.Valid {
		agentVersion = m.TargetAgentVersion.String
	}

	model.AgentVersion, err = version.Parse(agentVersion)
	if err != nil {
		return coremodel.ModelInfo{}, errors.Errorf(
			"parsing model %q agent version %q: %w",
			m.UUID, agentVersion, err,
		)
	}

	model.ControllerUUID, err = uuid.UUIDFromString(m.ControllerUUID)
	if err != nil {
		return coremodel.ModelInfo{}, errors.Errorf(
			"parsing controller uuid %q for model %q: %w",
			m.ControllerUUID, m.UUID, err,
		)
	}
	return model, nil
}

// GetModelMetrics the current model info and its associated metrics.
// If no model has been set then an error satisfying
// [modelerrors.NotFound] is returned.
func (s *ModelState) GetModelMetrics(ctx context.Context) (coremodel.ModelMetrics, error) {
	readOnlyModel, err := s.GetModel(ctx)
	if err != nil {
		return coremodel.ModelMetrics{}, err
	}

	db, err := s.DB()
	if err != nil {
		return coremodel.ModelMetrics{}, errors.Capture(err)
	}

	var modelMetrics dbModelMetrics
	stmt, err := s.Prepare(`SELECT &dbModelMetrics.* FROM v_model_metrics;`, modelMetrics)
	if err != nil {
		return coremodel.ModelMetrics{}, errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt).Get(&modelMetrics)
		if err != nil {
			return errors.Errorf("getting model metrics: %w", err)
		}
		return nil
	})
	if err != nil {
		return coremodel.ModelMetrics{}, err
	}

	return coremodel.ModelMetrics{
		Model:            readOnlyModel,
		ApplicationCount: modelMetrics.ApplicationCount,
		MachineCount:     modelMetrics.MachineCount,
		UnitCount:        modelMetrics.UnitCount,
	}, nil
}

// GetModelCloudType returns the cloud type from a model that has been
// set in the database. If no model exists then an error satisfying
// [modelerrors.NotFound] is returned.
func (s *ModelState) GetModelCloudType(ctx context.Context) (string, error) {
	db, err := s.DB()
	if err != nil {
		return "", errors.Capture(err)
	}

	m := dbReadOnlyModel{}
	stmt, err := s.Prepare(`SELECT &dbReadOnlyModel.cloud_type FROM model`, m)
	if err != nil {
		return "", errors.Capture(err)
	}

	err = db.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		err := tx.Query(ctx, stmt).Get(&m)
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("model does not exist").Add(modelerrors.NotFound)
		}
		return err
	})

	if err != nil {
		return "", errors.Capture(err)
	}

	return m.CloudType, nil
}

// CreateReadOnlyModel is responsible for creating a new model within the model
// database. If the model already exists then an error satisfying
// [modelerrors.AlreadyExists] is returned.
func CreateReadOnlyModel(ctx context.Context, args model.ModelDetailArgs, preparer domain.Preparer, tx *sqlair.TX) error {
	// This is some defensive programming. The zero value of agent version is
	// still valid but should really be considered null for the purposes of
	// allowing the DDL to assert constraints.
	var agentVersion sql.NullString
	if args.AgentVersion != version.Zero {
		agentVersion.String = args.AgentVersion.String()
		agentVersion.Valid = true
	}

	uuid := dbUUID{UUID: args.UUID.String()}
	checkExistsStmt, err := preparer.Prepare(`
SELECT &dbUUID.uuid
FROM model
	`, uuid)
	if err != nil {
		return errors.Capture(err)
	}

	err = tx.Query(ctx, checkExistsStmt).Get(&uuid)
	if err != nil && !errors.Is(err, sqlair.ErrNoRows) {
		return errors.Errorf(
			"checking if model %q already exists: %w",
			args.UUID, err,
		)
	} else if err == nil {
		return errors.Errorf(
			"creating readonly model %q information but model already exists",
			args.UUID,
		).Add(modelerrors.AlreadyExists)
	}

	m := dbReadOnlyModel{
		UUID:               args.UUID.String(),
		ControllerUUID:     args.ControllerUUID.String(),
		Name:               args.Name,
		Type:               args.Type.String(),
		TargetAgentVersion: agentVersion,
		Cloud:              args.Cloud,
		CloudType:          args.CloudType,
		CloudRegion:        args.CloudRegion,
		CredentialOwner:    args.CredentialOwner.Name(),
		CredentialName:     args.CredentialName,
		IsControllerModel:  args.IsControllerModel,
	}

	insertStmt, err := preparer.Prepare(`
INSERT INTO model (*) VALUES ($dbReadOnlyModel.*)
`, dbReadOnlyModel{})
	if err != nil {
		return errors.Capture(err)
	}

	if err := tx.Query(ctx, insertStmt, m).Run(); err != nil {
		return errors.Errorf(
			"creating readonly model %q information: %w", args.UUID, err,
		)
	}

	return nil
}
