// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"time"

	"github.com/juju/juju/core/charm"
	"github.com/juju/juju/core/objectstore"
	coreunit "github.com/juju/juju/core/unit"
	"github.com/juju/juju/domain/application/architecture"
	domaincharm "github.com/juju/juju/domain/application/charm"
	"github.com/juju/juju/domain/ipaddress"
	"github.com/juju/juju/domain/linklayerdevice"
	internalcharm "github.com/juju/juju/internal/charm"
	charmresource "github.com/juju/juju/internal/charm/resource"
)

// AddApplicationArg contains parameters for saving an application to state.
type AddApplicationArg struct {
	Charm             domaincharm.Charm
	CharmDownloadInfo *domaincharm.DownloadInfo
	Scale             int
	Platform          Platform
	Channel           *Channel
	// Resources defines the list of resources to add to an application.
	// They should match all the resources defined in the Charm.
	Resources []AddApplicationResourceArg
}

// AddApplicationResourceArg defines the arguments required to add a resource to an application.
type AddApplicationResourceArg struct {
	Name     string
	Revision *int
	Origin   charmresource.Origin
}

// Channel represents the channel of a application charm.
// Do not confuse this with a channel that is in the manifest file found
// in the charm package. They represent different concepts, yet hold the
// same data.
type Channel struct {
	Track  string
	Risk   ChannelRisk
	Branch string
}

// ChannelRisk describes the type of risk in a current channel.
type ChannelRisk string

const (
	RiskStable    ChannelRisk = "stable"
	RiskCandidate ChannelRisk = "candidate"
	RiskBeta      ChannelRisk = "beta"
	RiskEdge      ChannelRisk = "edge"
)

// OSType represents the type of an application's OS.
type OSType int

const (
	Ubuntu OSType = iota
)

// Platform contains parameters for an application's platform.
type Platform struct {
	Channel      string
	OSType       OSType
	Architecture Architecture
}

// Architecture represents the architecture of a application charm.
type Architecture = architecture.Architecture

// ScaleState describes the scale status of a k8s application.
type ScaleState struct {
	Scaling     bool
	Scale       int
	ScaleTarget int
}

// CloudService contains parameters for an application's cloud service.
type CloudService struct {
	ProviderId string
	Address    *ServiceAddress
}

// ServiceAddress contains parameters for a cloud service address.
// This may be from a load balancer, or cluster service etc.
type ServiceAddress struct {
	Value       string
	AddressType ipaddress.AddressType
	Scope       ipaddress.Scope
	Origin      ipaddress.Origin
	ConfigType  ipaddress.ConfigType
}

// Origin contains parameters for an application's origin.
type Origin struct {
	ID       string
	Channel  Channel
	Platform Platform
	Revision int
}

const (
	// HashAlgorithmSHA256 is the sha256 hash algorithm.
	// Currently it's the only one.
	HashAlgorithmSHA256 = 0
)

// PasswordInfo contains password parameters.
type PasswordInfo struct {
	PasswordHash  string
	HashAlgorithm int
}

// CloudContainer contains parameters for a unit's cloud container.
type CloudContainer struct {
	ProviderId string
	Address    *ContainerAddress
	Ports      *[]string
}

// ContainerDevice is the placeholder link layer device
// used to tie the cloud container IP address to the container.
type ContainerDevice struct {
	Name              string
	DeviceTypeID      linklayerdevice.DeviceType
	VirtualPortTypeID linklayerdevice.VirtualPortType
}

// ContainerAddress contains parameters for a cloud container address.
// Device is an attribute of address rather than cloud container
// since it's a placeholder used to tie the address to the
// cloud container and is only needed if the address exists.
type ContainerAddress struct {
	Device      ContainerDevice
	Value       string
	AddressType ipaddress.AddressType
	Scope       ipaddress.Scope
	Origin      ipaddress.Origin
	ConfigType  ipaddress.ConfigType
}

// AddUnitArg contains parameters for adding a unit to state.
type AddUnitArg struct {
	UnitName coreunit.Name
	UnitStatusArg
}

// InsertUnitArg is used to insert a fully populated unit.
// Used by import and when registering a CAAS unit.
type InsertUnitArg struct {
	UnitName       coreunit.Name
	CloudContainer *CloudContainer
	Password       *PasswordInfo
	UnitStatusArg
}

// UnitStatusArg contains parameters for updating a unit status in state.
type UnitStatusArg struct {
	AgentStatus    UnitAgentStatusInfo
	WorkloadStatus UnitWorkloadStatusInfo
}

// StatusInfo holds details about the status of an entity.
type StatusInfo struct {
	Message string
	Data    map[string]string
	Since   time.Time
}

// CloudContainerStatusStatusInfo holds a cloud container status
// and associated information.
type CloudContainerStatusStatusInfo struct {
	StatusID CloudContainerStatusType
	StatusInfo
}

// UnitAgentStatusInfo holds a unit agent status
// and associated information.
type UnitAgentStatusInfo struct {
	StatusID UnitAgentStatusType
	StatusInfo
}

// UnitWorkloadStatusInfo holds a unit workload status
// and associated information.
type UnitWorkloadStatusInfo struct {
	StatusID UnitWorkloadStatusType
	StatusInfo
}

// CharmDownloadInfo contains parameters for downloading a charm.
type CharmDownloadInfo struct {
	CharmUUID    charm.ID
	Name         string
	SHA256       string
	DownloadInfo domaincharm.DownloadInfo
}

// ResolveCharmDownload contains parameters for resolving a charm download.
type ResolveCharmDownload struct {
	CharmUUID charm.ID
	SHA256    string
	SHA384    string
	Path      string
	Size      int64
}

// ResolveControllerCharmDownload contains parameters for resolving a charm
// download.
type ResolveControllerCharmDownload struct {
	SHA256 string
	SHA384 string
	Path   string
	Size   int64
}

// ResolvedCharmDownload contains parameters for a resolved charm download.
type ResolvedCharmDownload struct {
	// Actions is the actions that the charm supports.
	// Deprecated: should be filled in by the charm store.
	Actions         domaincharm.Actions
	LXDProfile      []byte
	ObjectStoreUUID objectstore.UUID
	ArchivePath     string
}

// ResolvedControllerCharmDownload contains parameters for a resolved controller
// charm download.
type ResolvedControllerCharmDownload struct {
	Charm           internalcharm.Charm
	ArchivePath     string
	ObjectStoreUUID objectstore.UUID
}

// RevisionUpdaterApplication is responsible for updating the revision of an
// application.
type RevisionUpdaterApplication struct {
	Name         string
	CharmLocator domaincharm.CharmLocator
	Origin       Origin
	NumUnits     int
}
