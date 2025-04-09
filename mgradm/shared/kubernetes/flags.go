// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
)

// KubernetesServerFlags is the aggregation of all flags for install, upgrade and migrate.
type KubernetesServerFlags struct {
	Mirror            string
	TZ                string
	utils.ServerFlags `mapstructure:",squash"`
	Kubernetes        utils.KubernetesFlags
	Volumes           utils.VolumesFlags
	// SSH defines the SSH configuration to use to connect to the source server to migrate.
	SSH          utils.SSHFlags
	Installation utils.InstallationFlags `mapstructure:",squash"`
	Migration    utils.MigrationFlags    `mapstructure:",squash"`
	Upgrade      utils.UpgradeFlags      `mapstructure:",squash"`
}
