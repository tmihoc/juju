// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"time"

	coreapplication "github.com/juju/juju/core/application"
	coreresource "github.com/juju/juju/core/resource"
	"github.com/juju/juju/domain/resource"
	charmresource "github.com/juju/juju/internal/charm/resource"
	"github.com/juju/juju/internal/errors"
)

// resourceIdentity represents the unique identity of a resource within an
// application.
type resourceIdentity struct {
	UUID            string `db:"uuid"`
	ApplicationUUID string `db:"application_uuid"`
	Name            string `db:"name"`
}

// resourceUUID represents the unique identifier of a resource.
type resourceUUID struct {
	UUID string `db:"uuid"`
}

// resourceKind is the kind of the resource, e.g. file or oci-image.
type resourceKind struct {
	Name string `db:"kind_name"`
	UUID string `db:"uuid"`
}

// resourceView represents the view model for a resource entity. It contains
// all fields from v_resource table view.
type resourceView struct {
	UUID            string    `db:"uuid"`
	ApplicationUUID string    `db:"application_uuid"`
	Name            string    `db:"name"`
	CreatedAt       time.Time `db:"created_at"`
	Revision        int       `db:"revision"`
	OriginTypeId    int       `db:"origin_type_id"`
	RetrievedBy     string    `db:"retrieved_by"`
	RetrievedByType string    `db:"retrieved_by_type"`
	Path            string    `db:"path"`
	Description     string    `db:"description"`
	Kind            string    `db:"kind_name"`
	Size            int64     `db:"size"`
	SHA384          string    `db:"sha384"`
}

// toCharmResource converts the resourceView struct to a
// charmresource.Resource, populating its fields accordingly.
func (rv resourceView) toCharmResource() (charmresource.Resource, error) {
	kind, err := charmresource.ParseType(rv.Kind)
	if err != nil {
		return charmresource.Resource{}, errors.Errorf("converting resource type: %w", err)
	}
	var fingerprint charmresource.Fingerprint
	if rv.SHA384 != "" {
		fingerprint, err = charmresource.ParseFingerprint(rv.SHA384)
		if err != nil {
			return charmresource.Resource{}, errors.Errorf("converting resource fingerprint: %w", err)
		}
	}

	return charmresource.Resource{
		Meta: charmresource.Meta{
			Name:        rv.Name,
			Type:        kind,
			Path:        rv.Path,
			Description: rv.Description,
		},
		Origin:      charmresource.Origin(rv.OriginTypeId),
		Revision:    rv.Revision,
		Fingerprint: fingerprint,
		Size:        rv.Size,
	}, nil
}

// toResource converts a resourceView object to a resource.Resource object
// including metadata and timestamps.
func (rv resourceView) toResource() (resource.Resource, error) {
	charmRes, err := rv.toCharmResource()
	if err != nil {
		return resource.Resource{}, errors.Capture(err)
	}
	return resource.Resource{
		Resource:        charmRes,
		UUID:            coreresource.UUID(rv.UUID),
		ApplicationID:   coreapplication.ID(rv.ApplicationUUID),
		RetrievedBy:     rv.RetrievedBy,
		RetrievedByType: resource.RetrievedByType(rv.RetrievedByType),
		Timestamp:       rv.CreatedAt,
	}, nil
}

// unitResource represents the mapping of a resource to a unit.
type unitResource struct {
	ResourceUUID string    `db:"resource_uuid"`
	UnitUUID     string    `db:"unit_uuid"`
	AddedAt      time.Time `db:"added_at"`
}

type applicationNameAndID struct {
	ApplicationID coreapplication.ID `db:"uuid"`
	Name          string             `db:"name"`
}

// kubernetesApplicationResource represents the mapping of a resource to a unit.
type kubernetesApplicationResource struct {
	ResourceUUID string    `db:"resource_uuid"`
	AddedAt      time.Time `db:"added_at"`
}

type storedFileResource struct {
	ObjectStoreUUID string `db:"store_uuid"`
	ResourceUUID    string `db:"resource_uuid"`
	Size            int64  `db:"size"`
	SHA384          string `db:"sha384"`
}

type storedContainerImageResource struct {
	StorageKey   string `db:"store_storage_key"`
	ResourceUUID string `db:"resource_uuid"`
	Size         int64  `db:"size"`
	Hash         string `db:"sha384"`
}
