// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd shared_podman.Systemd = shared_podman.SystemdImpl{}

func upgradePodman(_ *types.GlobalFlags, flags *podmanUpgradeFlags, cmd *cobra.Command, _ []string) error {
	hostData, err := shared_podman.InspectHost()
	if err != nil {
		return err
	}

	flags.ServerFlags.CheckParameters()
	flags.Installation.CheckUpgradeParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	authFile, cleaner, err := shared_podman.PodmanLogin(hostData, flags.Installation.SCC, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to login to %s"), flags.Image.RegistryFQDN)
	}
	defer cleaner()

	return podman.Upgrade(
		systemd, authFile,
		flags.Installation.DB,
		flags.Installation.ReportDB,
		flags.Installation.SSL,
		flags.Image,
		flags.DBUpgradeImage,
		flags.Coco,
		flags.HubXmlrpc,
		flags.Saline,
		flags.Pgsql,
		flags.Installation.SCC,
		flags.Installation.TZ,
	)
}
