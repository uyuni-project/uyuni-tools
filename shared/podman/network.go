// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// The name of the podman network for Uyuni and its proxies.
const UyuniNetwork = "uyuni"

func hasIpv6Enabled(network string) bool {
	hasIpv6, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "network", "inspect",
		"--format", "{{.IPv6Enabled}}", network)
	if err == nil && strings.TrimSpace(string(hasIpv6)) == "true" {
		return true
	}
	return false
}

// SetupNetwork creates the podman network.
func SetupNetwork(isProxy bool) error {
	log.Info().Msgf(L("Setting up %s network"), UyuniNetwork)

	ipv6Enabled := isIpv6Enabled()

	// check if network exists before trying to get the IPV6 information
	networkExists := IsNetworkPresent(UyuniNetwork)
	if networkExists {
		log.Debug().Msgf("%s network already present", UyuniNetwork)
		// Check if the uyuni network exists and is IPv6 enabled
		hasIpv6 := hasIpv6Enabled(UyuniNetwork)
		if !hasIpv6 && ipv6Enabled {
			log.Info().Msgf(L("%s network doesn't have IPv6, deleting existing network to enable IPv6 on it"), UyuniNetwork)
			err := utils.RunCmd("podman", "network", "rm", UyuniNetwork,
				"--log-level", log.Logger.GetLevel().String())
			if err != nil {
				return utils.Errorf(err, L("failed to remove %s podman network"), UyuniNetwork)
			}
		} else {
			log.Info().Msgf(L("Reusing existing %s network"), UyuniNetwork)
			return nil
		}
	}

	// We do not need inter-container resolution, disable dns plugin
	args := []string{"network", "create"}
	if isProxy {
		args = append(args, "--disable-dns")
	}
	if ipv6Enabled {
		// An IPv6 network on a host where IPv6 is disabled doesn't work: don't try it.
		// Check if the networkd backend is netavark
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "info", "--format", "{{.Host.NetworkBackend}}")
		backend := strings.Trim(string(out), "\n")
		if err != nil {
			return utils.Errorf(err, L("failed to find podman's network backend"))
		} else if backend != "netavark" {
			log.Info().Msgf(L("Podman's network backend (%[1]s) is not netavark, skipping IPv6 enabling on %[2]s network"),
				backend, UyuniNetwork)
		} else {
			args = append(args, "--ipv6")
		}
	}
	args = append(args, UyuniNetwork)
	err := utils.RunCmd("podman", args...)
	if err != nil {
		return utils.Errorf(err, L("failed to create %s network with IPv6 enabled"), UyuniNetwork)
	}
	return nil
}

func isIpv6Enabled() bool {
	files := []string{
		"/sys/module/ipv6/parameters/disable",
		"/proc/sys/net/ipv6/conf/default/disable_ipv6",
		"/proc/sys/net/ipv6/conf/all/disable_ipv6",
	}

	for _, file := range files {
		// Mind that we are checking disable files, the semantic is inverted
		if utils.GetFileBoolean(file) {
			return false
		}
	}
	return true
}

// DeleteNetwork deletes the uyuni podman network.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteNetwork(dryRun bool) {
	err := utils.RunCmd("podman", "network", "exists", UyuniNetwork)
	if err != nil {
		log.Info().Msgf(L("Network %s already removed"), UyuniNetwork)
	} else {
		if dryRun {
			log.Info().Msgf(L("Would run %s"), "podman network rm "+UyuniNetwork)
		} else {
			err := utils.RunCmd("podman", "network", "rm", UyuniNetwork)
			if err != nil {
				log.Error().Msgf(L("Failed to remove network %s"), UyuniNetwork)
			} else {
				log.Info().Msg(L("Network removed"))
			}
		}
	}
}

// IsNetworkPresent returns whether a network is already present.
func IsNetworkPresent(network string) bool {
	cmd := exec.Command("podman", "network", "exists", network)
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}
