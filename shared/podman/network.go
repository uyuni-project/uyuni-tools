// SPDX-FileCopyrightText: 2023 SUSE LLC
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

// The name of the podman network for Uyuni and its proxies
const UyuniNetwork = "uyuni"

func SetupNetwork() {
	log.Info().Msgf("Setting up %s network", UyuniNetwork)

	ipv6Enabled := isIpv6Enabled()

	testNetworkCmd := exec.Command("podman", "network", "exists", UyuniNetwork)
	testNetworkCmd.Run()
	// check if network exists before trying to get the IPV6 information
	if testNetworkCmd.ProcessState.ExitCode() == 0 {
		// Check if the uyuni network exists and is IPv6 enabled
		hasIpv6, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "network", "inspect", "--format", "{{.IPv6Enabled}}", UyuniNetwork)
		if err == nil {
			if string(hasIpv6) != "true" && ipv6Enabled {
				log.Info().Msgf("%s network doesn't have IPv6, deleting existing network to enable IPv6 on it", UyuniNetwork)
				message := fmt.Sprintf("Failed to remove %s podman network", UyuniNetwork)
				err := utils.RunCmd("podman", "network", "rm", UyuniNetwork,
					"--log-level", log.Logger.GetLevel().String())
				if err != nil {
					log.Fatal().Err(err).Msg(message)
				}
			} else {
				log.Info().Msgf("Reusing existing %s network", UyuniNetwork)
				return
			}
		}
	}

	message := fmt.Sprintf("Failed to create %s network with IPv6 enabled", UyuniNetwork)

	args := []string{"network", "create"}
	if ipv6Enabled {
		// An IPv6 network on a host where IPv6 is disabled doesn't work: don't try it.
		// Check if the networkd backend is netavark
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "info", "--format", "{{.Host.NetworkBackend}}")
		backend := strings.Trim(string(out), "\n")
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to find podman's network backend")
		} else if backend != "netavark" {
			log.Info().Msgf("Podman's network backend (%s) is not netavark, skipping IPv6 enabling on %s network", backend, UyuniNetwork)
		} else {
			args = append(args, "--ipv6")
		}
	}
	args = append(args, UyuniNetwork)
	err := utils.RunCmd("podman", args...)
	if err != nil {
		log.Fatal().Err(err).Msg(message)
	}
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
