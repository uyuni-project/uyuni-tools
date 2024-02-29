// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// The name of the podman network for Uyuni and its proxies.
const UyuniNetwork = "uyuni"

// SetupNetwork creates the podman network.
func SetupNetwork() error {
	log.Info().Msgf("Setting up %s network", UyuniNetwork)

	ipv6Enabled := isIpv6Enabled()

	// check if network exists before trying to get the IPV6 information
	networkExists := IsNetworkPresent(UyuniNetwork)
	if networkExists {
		log.Debug().Msgf("%s network already present", UyuniNetwork)
		// Check if the uyuni network exists and is IPv6 enabled
		hasIpv6, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "network", "inspect", "--format", "{{.IPv6Enabled}}", UyuniNetwork)
		if err == nil {
			if string(hasIpv6) != "true" && ipv6Enabled {
				log.Info().Msgf("%s network doesn't have IPv6, deleting existing network to enable IPv6 on it", UyuniNetwork)
				err := utils.RunCmd("podman", "network", "rm", UyuniNetwork,
					"--log-level", log.Logger.GetLevel().String())
				if err != nil {
					return fmt.Errorf("failed to remove %s podman network: %s", UyuniNetwork, err)
				}
			} else {
				log.Info().Msgf("Reusing existing %s network", UyuniNetwork)
				return nil
			}
		}
	}

	args := []string{"network", "create"}
	if ipv6Enabled {
		// An IPv6 network on a host where IPv6 is disabled doesn't work: don't try it.
		// Check if the networkd backend is netavark
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "info", "--format", "{{.Host.NetworkBackend}}")
		backend := strings.Trim(string(out), "\n")
		if err != nil {
			return fmt.Errorf("failed to find podman's network backend: %s", err)
		} else if backend != "netavark" {
			log.Info().Msgf("Podman's network backend (%s) is not netavark, skipping IPv6 enabling on %s network", backend, UyuniNetwork)
		} else {
			args = append(args, "--ipv6")
		}
	}
	args = append(args, UyuniNetwork)
	err := utils.RunCmd("podman", args...)
	if err != nil {
		return fmt.Errorf("failed to create %s network with IPv6 enabled: %s", UyuniNetwork, err)
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
		log.Info().Msgf("Network %s already removed", UyuniNetwork)
	} else {
		if dryRun {
			log.Info().Msgf("Would run podman network rm %s", UyuniNetwork)
		} else {
			err := utils.RunCmd("podman", "network", "rm", UyuniNetwork)
			if err != nil {
				log.Error().Msgf("Failed to remove network %s", UyuniNetwork)
			} else {
				log.Info().Msg("Network removed")
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
