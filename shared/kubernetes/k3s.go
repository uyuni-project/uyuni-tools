// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const k3sTraefikConfigPath = "/var/lib/rancher/k3s/server/manifests/uyuni-traefik-config.yaml"
const k3sTraefikMainConfigPath = "/var/lib/rancher/k3s/server/manifests/traefik.yaml"

// InstallK3sTraefikConfig install K3s Traefik configuration.
func InstallK3sTraefikConfig(ports []types.PortMap) error {
	log.Info().Msg(L("Installing K3s Traefik configuration"))

	endpoints := []types.PortMap{}
	for _, port := range ports {
		port.Name = getTraefikEndpointName(port)
		endpoints = append(endpoints, port)
	}
	version, err := getTraefikChartMajorVersion()
	if err != nil {
		return err
	}

	data := K3sTraefikConfigTemplateData{
		Ports:         endpoints,
		ExposeBoolean: version < 27,
	}
	if err := utils.WriteTemplateToFile(data, k3sTraefikConfigPath, 0600, true); err != nil {
		return utils.Errorf(err, L("Failed to write Traefik configuration"))
	}

	// Wait for traefik to be back
	return waitForTraefik()
}

// getTraefikEndpointName computes the traefik endpoint name from the service and port names.
// Those names should be less than 15 characters long.
func getTraefikEndpointName(portmap types.PortMap) string {
	svc := shortenName(portmap.Service)
	name := shortenName(portmap.Name)
	if name != svc {
		return fmt.Sprintf("%s-%s", svc, name)
	}
	return name
}

func shortenName(name string) string {
	shorteningMap := map[string]string{
		"taskomatic":      "tasko",
		"metrics":         "mtrx",
		"postgresql":      "pgsql",
		"exporter":        "xport",
		"uyuni-proxy-tcp": "uyuni",
		"uyuni-proxy-udp": "uyuni",
	}
	short := shorteningMap[name]
	if short == "" {
		short = name
	}
	return short
}

var newRunner = utils.NewRunner

func waitForTraefik() error {
	log.Info().Msg(L("Waiting for Traefik to be reloaded"))
	for i := 0; i < 120; i++ {
		out, err := newRunner("kubectl", "get", "job", "-n", "kube-system",
			"-o", "jsonpath={.status.completionTime}", "helm-install-traefik").Log(zerolog.TraceLevel).Exec()
		if err == nil {
			completionTime, err := time.Parse(time.RFC3339, string(out))
			if err == nil && time.Since(completionTime.Local()).Seconds() < 60 {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New(L("Failed to reload Traefik"))
}

// UninstallK3sTraefikConfig uninstall K3s Traefik configuration.
func UninstallK3sTraefikConfig(dryRun bool) {
	// Write a blank file first to get traefik to be reinstalled
	if !dryRun {
		log.Info().Msg(L("Reinstalling Traefik without additionnal configuration"))
		err := os.WriteFile(k3sTraefikConfigPath, []byte{}, 0600)
		if err != nil {
			log.Error().Err(err).Msg(L("failed to write empty traefik configuration"))
		} else {
			// Wait for traefik to be back
			if err := waitForTraefik(); err != nil {
				log.Error().Err(err).Msg(L("failed to uninstall traefik configuration"))
			}
		}
	} else {
		log.Info().Msg(L("Would reinstall Traefik without additionnal configuration"))
	}

	// Now that it's reinstalled, remove the file
	utils.UninstallFile(k3sTraefikConfigPath, dryRun)
}

func getTraefikChartMajorVersion() (int, error) {
	out, err := os.ReadFile(k3sTraefikMainConfigPath)
	if err != nil {
		return 0, utils.Errorf(err, L("failed to read the traefik configuration"))
	}
	matches := regexp.MustCompile(`traefik-([0-9]+)`).FindStringSubmatch(string(out))
	if matches == nil {
		return 0, errors.New(L("traefik configuration file doesn't contain the helm chart version"))
	}
	if len(matches) != 2 {
		return 0, errors.New(L("failed to find traefik helm chart version"))
	}

	majorVersion, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, utils.Errorf(err, L(""))
	}

	return majorVersion, nil
}
