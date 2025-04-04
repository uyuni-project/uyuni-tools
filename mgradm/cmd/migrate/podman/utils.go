// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	podman_utils "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman_utils.Systemd = podman_utils.SystemdImpl{}

func migrateToPodman(
	_ *types.GlobalFlags,
	flags *podmanMigrateFlags,
	cmd *cobra.Command,
	args []string,
) error {
	hostData, err := podman_utils.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman_utils.PodmanLogin(hostData, flags.SCC)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	flags.UpgradeFlags.CheckParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	return podman.Migrate(
		systemd, authFile,
		flags.Image.Registry,
		flags.ServerFlags.DB,
		flags.ReportDB,
		flags.SSL,
		flags.Image,
		flags.DBUpgradeImage,
		flags.Coco,
		flags.HubXmlrpc,
		flags.Saline,
		flags.Pgsql,
		flags.SCC,
		flags.TZ,
		flags.Prepare,
		flags.User,
		flags.Debug.Java,
		flags.Mirror,
		flags.Podman,
		args,
	)
}
