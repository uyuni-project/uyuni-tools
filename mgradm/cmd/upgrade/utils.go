// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package upgrade

import (
	"errors"
	"os"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd shared_podman.Systemd = shared_podman.NewSystemd()

func checkAndCleanIPv6PortBindings() error {
	// Workaround for IPv6 options that do not work with nftables - bsc#1268755
	servicePath := shared_podman.GetServicePath(shared_podman.ServerService)
	if !utils.FileExists(servicePath) {
		return nil
	}

	content, err := os.ReadFile(servicePath)
	if err != nil {
		return utils.Errorf(err, L("failed to read %s"), servicePath)
	}

	// For example "-p [::]:443:443 \"
	re := regexp.MustCompile(`(?m)^[ \t]*-p[ \t]+["']?\[::\]:\d+:\d+(?:/\w+)?["']?[ \t]*\\?[ \t]*(?:\r?\n|$)`)

	if !re.Match(content) {
		return nil
	}

	cleaned := re.ReplaceAllString(string(content), "")

	if err := os.Chmod(servicePath, 0640); err != nil {
		return utils.Errorf(err, L("failed to change permissions for %s"), servicePath)
	}

	if err := os.WriteFile(servicePath, []byte(cleaned), 0640); err != nil {
		return utils.Errorf(err, L("failed to write %s"), servicePath)
	}

	return errors.New(L("Legacy dual-stack IPv6 port bindings have been removed from uyuni-server.service. " +
		"Please reboot the host and then run the upgrade command again"))
}

func upgradePodman(_ *types.GlobalFlags, flags *podmanUpgradeFlags, cmd *cobra.Command, _ []string) error {
	if err := checkAndCleanIPv6PortBindings(); err != nil {
		return err
	}

	hostData, err := shared_podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := shared_podman.PodmanLogin(hostData, flags.Image.Registry, flags.Installation.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	flags.Installation.CheckUpgradeParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}
	podman.WarnIfServicesDisabled(systemd)

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
		flags.TFTPD,
		flags.Installation.TZ,
		flags.Installation.Debug.Java,
	)
}
