// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Server

// UyuniServices is the list of services to expose.
var UyuniServices = []types.UyuniService{
	{Name: "uyuni-server",
		Image:       ServerImage,
		Description: L("Main service"),
		Replicas:    types.SingleMandatoryReplica,
		Options:     []types.UyuniServiceOption{}},

	{Name: "uyuni-server-migration",
		Image:       Migration14To16Image,
		Description: L("Migration helper"),
		Replicas:    types.SingleOptionalReplica,
		Options:     []types.UyuniServiceOption{}},

	{Name: "uyuni-server-attestation",
		Image:       COCOAttestationImage,
		Description: L("Confidential computing attestation"),
		Replicas:    types.SingleOptionalReplica,
		Options:     []types.UyuniServiceOption{}},

	{Name: "uyuni-hub-xmlrpc",
		Image:       HubXMLRPCImage,
		Description: L("Hub XML-RPC API"),
		Replicas:    types.SingleOptionalReplica,
		Options:     []types.UyuniServiceOption{}},

	{Name: "uyuni-saline",
		Image:       SalineImage,
		Description: L("Saline"),
		Replicas:    types.SingleOptionalReplica,
		Options:     []types.UyuniServiceOption{}},
	{Name: "uyuni-db",
		Description: L("Database"),
		Replicas:    types.SingleMandatoryReplica,
		Options:     []types.UyuniServiceOption{}},
}

// ServerImage holds the flags to tune the server container image.
var ServerImage = types.ImageFlags{
	Name:       "server",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

// HubXMLRPCImage holds the flags to tune the hub XML-RPC API container image.
var HubXMLRPCImage = types.ImageFlags{
	Name:       "server-hub-xmlrpc-api",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

// COCOAttestationImage holds the flags to tune the confidential computing attestation container image.
var COCOAttestationImage = types.ImageFlags{
	Name:       "server-attestation",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

// Saline holds the flags to tune the saline container image.
var SalineImage = types.ImageFlags{
	Name:       "server-saline",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

// Migration14To16Image holds the flags to tune the DB migration container image.
var Migration14To16Image = types.ImageFlags{
	Name:       "server-migration-14-16",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

// PostgreSQLImage holds the flags to tune the DB container image.
var PostgreSQLImage = types.ImageFlags{
	Name:       "uyuni-db",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}
