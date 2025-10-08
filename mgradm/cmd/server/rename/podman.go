// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package rename

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_podman "github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.NewSystemd()

func renameForPodman(_ *types.GlobalFlags, flags *renameFlags, _ *cobra.Command, args []string) error {
	fqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}

	// Regenerate Server SSL certificate if requested
	image := podman.GetServiceImage(podman.ServerService)
	tz := findTimezone("podman")
	log.Info().Msg(L("Preparing SSL certificates to match the new hostname"))
	if err := adm_podman.PrepareSSLCertificates(image, &flags.SSL, tz, fqdn); err != nil {
		return err
	}

	log.Info().Msg(L("Stopping the server container"))
	if err := systemd.StopService(podman.ServerService); err != nil {
		return err
	}

	log.Info().Msgf(L("Changing the UYUNI_HOSTNAME to %s"), fqdn)
	if err := alterHostnameConfig(fqdn); err != nil {
		return err
	}

	// Update the service to ensure it has -e UYUNI_HOSTNAME
	if err := adm_podman.UpdateServerSystemdService(); err != nil {
		return err
	}

	// Restart the server container: the UYUNI_HOSTNAME change will be picked up by the uyuni-update-config service
	log.Info().Msg(L("Starting the server container"))
	err = systemd.StartService(podman.ServerService)
	if err != nil {
		return err
	}

	log.Info().Msg(L(`The renaming continues inside the server container.
The logs can be found in journalctl -u uyuni-config-update.service output.`))
	return nil
}

func findTimezone(backend string) string {
	// If the container is running, call podman exec uyuni-server sh -c "echo $TZ".
	cnx := shared.NewConnection(backend, podman.ServerContainerName, "")
	out, err := cnx.Exec("echo", "$TZ")
	if err == nil {
		return strings.TrimSpace(string(out))
	}

	if backend == "podman" {
		// Otherwise get the value from the uyuni-server.service.d/custom.conf file, in the 'Environment=TZ=' line.
		// In theory users shouldn't remove this, but who knows what could happen?
		log.Debug().Msg("Failed to get the timezone from the container, looking for it in systemd configuration file")
		if env, err := systemd.Show(podman.ServerService, "Environment"); err == nil {
			pattern := regexp.MustCompile("TZ=([^[:space]]*)")
			matches := pattern.FindStringSubmatch(env)
			if len(matches) == 1 {
				return matches[0]
			}
		}
	}

	log.Debug().Msg("Failed to get the timezone from the configuration, getting the host one")
	return utils.GetLocalTimezone()
}

// alterHostnameConfig changes the UYUNI_HOSTNAME value in the server systemd service or adds it if needed.
func alterHostnameConfig(fqdn string) error {
	config, err := readCustomConf()
	if err != nil {
		return err
	}

	// Append Environment=UYUNI_HOSTNAME={{.fqdn}} or replace the value
	pattern := regexp.MustCompile(`(?m)^Environment=UYUNI_HOSTNAME=.*$`)
	newConfig := pattern.ReplaceAllString(config, "Environment=UYUNI_HOSTNAME="+fqdn)
	if config == newConfig {
		newConfig = fmt.Sprintf("%s\nEnvironment=UYUNI_HOSTNAME=%s\n", config, fqdn)
	}

	systemdConfPath := podman.GetServiceConfFolder(podman.ServerService)
	customConf := path.Join(systemdConfPath, podman.CustomConf)
	if err := os.WriteFile(customConf, []byte(newConfig), 0640); err != nil {
		return utils.Error(err, L("failed to write custom.conf with the new hostname"))
	}

	return systemd.ReloadDaemon(false)
}

func readCustomConf() (config string, err error) {
	systemdConfPath := podman.GetServiceConfFolder(podman.ServerService)
	customConf := path.Join(systemdConfPath, podman.CustomConf)
	out, err := os.ReadFile(customConf)
	if err != nil {
		return "", utils.Error(err, L("failed to read the custom.conf file"))
	}

	config = string(out)
	return config, nil
}
