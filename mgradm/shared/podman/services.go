// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	sharedPodman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var replicatedServerServices = []string{
	sharedPodman.ServerAttestationService,
	sharedPodman.HubXmlrpcService,
	sharedPodman.SalineService,
}

var coreServerServices = []string{
	sharedPodman.ServerService,
	sharedPodman.DBService,
}

func enabledOptionalServices(systemd sharedPodman.Systemd) []string {
	services := []string{}
	if systemd.ServiceIsEnabled(sharedPodman.TFTPService) {
		services = append(services, sharedPodman.TFTPService)
	}
	for _, service := range replicatedServerServices {
		for i := 0; i < systemd.CurrentReplicaCount(service); i++ {
			services = append(services, fmt.Sprintf("%s@%d", service, i))
		}
	}
	return services
}

// ServerServicesEnabled reports whether the server is configured to start automatically at boot.
func ServerServicesEnabled(systemd sharedPodman.Systemd) bool {
	return systemd.ServiceIsEnabled(sharedPodman.ServerService)
}

// EnableServices enables automatic startup of the core server services without starting them.
// Optional services and replicas disabled by DisableServices need to be configured again explicitly.
func EnableServices(systemd sharedPodman.Systemd) error {
	if !systemd.HasService(sharedPodman.ServerService) {
		return errors.New(L("no installed server detected"))
	}

	errs := []error{}
	for _, service := range coreServerServices {
		if systemd.HasService(service) {
			errs = append(errs, systemd.EnableServiceAtBoot(service))
		}
	}
	return utils.JoinErrors(errs...)
}

// DisableServices disables automatic startup of all currently enabled server services without stopping them.
func DisableServices(systemd sharedPodman.Systemd) error {
	if !systemd.HasService(sharedPodman.ServerService) {
		return errors.New(L("no installed server detected"))
	}

	services := enabledOptionalServices(systemd)
	for _, service := range coreServerServices {
		if systemd.HasService(service) {
			services = append(services, service)
		}
	}

	errs := []error{}
	for _, service := range services {
		errs = append(errs, systemd.DisableServiceAtBoot(service))
	}
	return utils.JoinErrors(errs...)
}

// WarnIfServicesDisabled warns that the server stack will not start automatically at boot.
func WarnIfServicesDisabled(systemd sharedPodman.Systemd) {
	if !ServerServicesEnabled(systemd) {
		log.Warn().Msg(L("Server services are disabled and will not start automatically at boot"))
	}
}
