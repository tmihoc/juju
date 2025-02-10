// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package application

import (
	"time"
)

// StatusID represents the status of an entity.
type StatusID interface {
	UnsetStatusType | CloudContainerStatusType | UnitAgentStatusType | WorkloadStatusType
}

// StatusInfo holds details about the status of an entity.
type StatusInfo[T StatusID] struct {
	Status  T
	Message string
	Data    []byte
	Since   *time.Time
}

// UnsetStatusType represents the status of an entity that has not been set.
type UnsetStatusType int

const (
	UnsetStatus UnsetStatusType = iota
)

// CloudContainerStatusType represents the status of a cloud container
// as recorded in the cloud_container_status_value lookup table.
type CloudContainerStatusType int

const (
	CloudContainerStatusWaiting CloudContainerStatusType = iota
	CloudContainerStatusBlocked
	CloudContainerStatusRunning
)

// UnitAgentStatusType represents the status of a unit agent
// as recorded in the unit_agent_status_value lookup table.
type UnitAgentStatusType int

const (
	UnitAgentStatusAllocating UnitAgentStatusType = iota
	UnitAgentStatusExecuting
	UnitAgentStatusIdle
	UnitAgentStatusError
	UnitAgentStatusFailed
	UnitAgentStatusLost
	UnitAgentStatusRebooting
)

// WorkloadStatusType represents the status of a unit workload or application
// as recorded in the workload_status_value lookup table.
type WorkloadStatusType int

const (
	WorkloadStatusUnset WorkloadStatusType = iota
	WorkloadStatusUnknown
	WorkloadStatusMaintenance
	WorkloadStatusWaiting
	WorkloadStatusBlocked
	WorkloadStatusActive
	WorkloadStatusTerminated
)
