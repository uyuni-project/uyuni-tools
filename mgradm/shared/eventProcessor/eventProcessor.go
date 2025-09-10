// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package eventProcessor

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Upgrade event processor
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
	debugJava bool,
) error {
	if eventProcessorFlags.Image.Name == "" {
		return fmt.Errorf(L("image is required"))
	}

	if err := writeEventProcessorFiles(
		systemd, authFile, registry, eventProcessorFlags, baseImage, db, debugJava,
	); err != nil {
		return err
	}

	return systemd.ScaleService(1, podman.EventProcessorService)
}

func writeEventProcessorFiles(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
	debugJava bool,
) error {
	image := eventProcessorFlags.Image

	log.Debug().Msgf("Current running event processor replica is enforced to be 1")

	if !eventProcessorFlags.IsChanged {
		log.Debug().Msgf("Event processor settings are not changed.")
	}

	//eventProcessorImage, err := utils.ComputeImage(registry, baseImage.Tag, image)
	//if err != nil {
	//	return utils.Errorf(err, L("Failed to compute event processor image URL"))
	//}

	// TODO: Temporary solution to install image in my OBS branch, should be remove in production
	var eventProcessorImage string
	var err error

	// Check if the image name already contains a full registry path
	if strings.Contains(image.Name, "registry.opensuse.org") {
		eventProcessorImage = image.Name
		if image.Tag != "" && !strings.Contains(image.Name, ":") {
			eventProcessorImage += ":" + image.Tag
		} else if !strings.Contains(image.Name, ":") {
			eventProcessorImage += ":" + baseImage.Tag
		}
	} else {
		eventProcessorImage, err = utils.ComputeImage(registry, baseImage.Tag, image)
		if err != nil {
			return utils.Errorf(err, L("Failed to compute event processor image URL"))
		}
	}

	// Always enable pulling if service is requested (since we enforced single replica)
	preparedImage, err := podman.PrepareImage(authFile, eventProcessorImage, baseImage.PullPolicy, true)
	if err != nil {
		return err
	}

	// Control debug port expose from container to client
	var ports []types.PortMap
	if debugJava {
		ports = utils.EventProcessorPorts
	}

	eventProcessorData := templates.EventProcessorServiceTemplateData{
		NamePrefix:   "uyuni",
		Network:      podman.UyuniNetwork,
		DBUserSecret: podman.DBUserSecret,
		DBPassSecret: podman.DBPassSecret,
		DBBackend:    "postgres",
		Ports:        ports,
	}

	log.Info().Msg(L("Setting up event processor service"))

	if err := utils.WriteTemplateToFile(
		eventProcessorData, podman.GetServicePath(podman.EventProcessorService+"@"), 0555, true,
	); err != nil {
		return utils.Errorf(err, L("Failed to generate systemd service unit file"))
	}

	// Add conditional debug server inside container
	var javaOpts string
	if debugJava {
		javaOpts = "-Xdebug -Xrunjdwp:transport=dt_socket,address=*:8004,server=y,suspend=n"
	} else {
		javaOpts = ""
	}

	environment := fmt.Sprintf(`Environment=UYUNI_EVENT_PROCESSOR_IMAGE=%s
Environment=UYUNI_DB_NAME=%s
Environment=UYUNI_DB_PORT=%d
Environment=UYUNI_DB_HOST=%s
Environment=JAVA_OPTS="%s"`,
		preparedImage, db.Name, db.Port, db.Host, javaOpts)

	if err := podman.GenerateSystemdConfFile(
		podman.EventProcessorService+"@", "generated.conf", environment, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd user configuration file"))
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	return nil
}

func SetupEventProcessorContainer(
	systemd podman.Systemd,
	authFile string,
	registry string,
	eventProcessorFlags adm_utils.EventProcessorFlags,
	baseImage types.ImageFlags,
	db adm_utils.DBFlags,
	debugJava bool,
) error {
	if err := writeEventProcessorFiles(
		systemd, authFile, registry, eventProcessorFlags, baseImage, db, debugJava,
	); err != nil {
		return err
	}
	// Enforce one replica
	return systemd.ScaleService(1, podman.EventProcessorService)
}
