// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func GenerateSystemdService(httpdImage string, saltBrokerImage string, squidImage string, sshImage string,
	tftpdImage string, podmanArgs []string) {

	podman.SetupNetwork()

	log.Info().Msg("Generating systemd services")
	httpProxyConfig := getHttpProxyConfig()

	ports := []types.PortMap{}
	ports = append(ports, utils.PROXY_TCP_PORTS...)
	ports = append(ports, utils.UDP_PORTS...)

	// Pod
	dataPod := templates.PodTemplateData{
		Ports:         ports,
		HttpProxyFile: httpProxyConfig,
		Args:          strings.Join(podmanArgs, " "),
	}
	generateSystemdFile(dataPod, "pod")

	// Httpd
	dataHttpd := templates.HttpdTemplateData{
		Volumes:       utils.PROXY_HTTPD_VOLUMES,
		HttpProxyFile: httpProxyConfig,
		Image:         httpdImage,
	}
	generateSystemdFile(dataHttpd, "httpd")

	// Salt broker
	dataSaltBroker := templates.SaltBrokerTemplateData{
		HttpProxyFile: httpProxyConfig,
		Image:         saltBrokerImage,
	}
	generateSystemdFile(dataSaltBroker, "salt-broker")

	// Squid
	dataSquid := templates.SquidTemplateData{
		Volumes:       utils.PROXY_SQUID_VOLUMES,
		HttpProxyFile: httpProxyConfig,
		Image:         squidImage,
	}
	generateSystemdFile(dataSquid, "squid")

	// SSH
	dataSSH := templates.SSHTemplateData{
		HttpProxyFile: httpProxyConfig,
		Image:         sshImage,
	}
	generateSystemdFile(dataSSH, "ssh")

	// Tftpd
	dataTftpd := templates.TFTPDTemplateData{
		Volumes:       utils.PROXY_TFTPD_VOLUMES,
		HttpProxyFile: httpProxyConfig,
		Image:         tftpdImage,
	}
	generateSystemdFile(dataTftpd, "tftpd")

	utils.RunCmd("systemctl", "daemon-reload")
}

func generateSystemdFile(template utils.Template, service string) {
	name := fmt.Sprintf("uyuni-proxy-%s.service", service)

	const systemdPath = "/etc/systemd/system"
	path := path.Join(systemdPath, name)
	if err := utils.WriteTemplateToFile(template, path, 0555, true); err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate %s", path)
	}
}

func getHttpProxyConfig() string {
	const httpProxyConfigPath = "/etc/sysconfig/proxy"

	// Only SUSE distros seem to have such a file for HTTP proxy settings
	if utils.FileExists(httpProxyConfigPath) {
		return httpProxyConfigPath
	}
	return ""
}
