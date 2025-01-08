// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

// allowedEmbeddedCommands is a whitelist of Juju CLI commands which
// are permissible to run embedded on a controller.
var allowedEmbeddedCommands = []string{
	"actions",
	"add-machine",
	"add-space",
	"add-storage",
	"add-subnet",
	"add-unit",
	"add-user",
	"agreements",
	"attach",
	"attach-resource",
	"attach-storage",
	"bind",
	"cancel-task",
	"charm-resources",
	"clouds",
	"config",
	"consume",
	"controller-config",
	"create-storage-pool",
	"credentials",
	"deploy",
	"detach-storage",
	"disable-user",
	"enable-user",
	"expose",
	"find-offers",
	"firewall-rules",
	"constraints",
	"model-constraints",
	"help",
	"import-filesystem",
	"integrate",
	"machines",
	"metrics",
	"model-config",
	"model-default",
	"model-defaults",
	"move-to-space",
	"offer",
	"offers",
	"relate",
	"reload-spaces",
	"remove-application",
	"remove-credential",
	"remove-machine",
	"remove-offer",
	"remove-relation",
	"remove-saas",
	"remove-space",
	"remove-storage",
	"remove-storage-pool",
	"remove-unit",
	"remove-user",
	"rename-space",
	"resolved",
	"resolve",
	"resources",
	"resume-relation",
	"retry-provisioning",
	"run",
	"scale-application",
	"set-application-base",
	"set-constraints",
	"set-firewall-rule",
	"set-meter-status",
	"set-model-constraints",
	"show-action",
	"show-application",
	"show-cloud",
	"show-controller",
	"show-credential",
	"show-credentials",
	"show-machine",
	"show-model",
	"show-offer",
	"show-status-log",
	"show-storage",
	"show-space",
	"show-unit",
	"show-user",
	"sla",
	"spaces",
	"status",
	"storage",
	"storage-pools",
	"subnets",
	"suspend-relation",
	"trust",
	"unexpose",
	"update-storage-pool",
	"users",
}
