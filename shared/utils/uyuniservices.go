// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// Server

// Services
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

	{Name: "uyuni-server-saline",
		Image:       SalineImage,
		Description: L("Saline"),
		Replicas:    types.SingleOptionalReplica,
		Options:     []types.UyuniServiceOption{}},
	{Name: "uyuni-db",
		Description: L("Database"),
		Replicas:    types.SingleMandatoryReplica,
		Options:     []types.UyuniServiceOption{}},
}

// Images
var ServerImage = types.ImageFlags{
	Name:       "server",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

var HubXMLRPCImage = types.ImageFlags{
	Name:       "server-hub-xmlrpc-api",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

var COCOAttestationImage = types.ImageFlags{
	Name:       "server-attestation",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

var SalineImage = types.ImageFlags{
	Name:       "server-saline",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

var Migration14To16Image = types.ImageFlags{
	Name:       "server-migration-14-16",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}

var PostgreSQLImage = types.ImageFlags{
	Name:       "uyuni-db",
	Tag:        DefaultTag,
	Registry:   DefaultRegistry,
	PullPolicy: DefaultPullPolicy,
}
