// SPDX-FileCopyrightText: 2024 SUSE LLC
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

// GenerateSystemdService generates all the systemd files required by proxy.
func GenerateSystemdService(httpdImage string, saltBrokerImage string, squidImage string, sshImage string,
	tftpdImage string, podmanArgs []string) error {
	if err := podman.SetupNetwork(); err != nil {
		return fmt.Errorf("cannot setup network: %s", err)
	}

	log.Info().Msg("Generating systemd services")
	httpProxyConfig := getHttpProxyConfig()

	ports := []types.PortMap{}
	ports = append(ports, utils.PROXY_TCP_PORTS...)
	ports = append(ports, utils.PROXY_PODMAN_PORTS...)
	ports = append(ports, utils.UDP_PORTS...)

	// Pod
	dataPod := templates.PodTemplateData{
		Ports:         ports,
		HttpProxyFile: httpProxyConfig,
		Args:          strings.Join(podmanArgs, " "),
	}
	if err := generateSystemdFile(dataPod, "pod"); err != nil {
		return fmt.Errorf("cannot generated systemd file: %s", err)
	}

	// Httpd
	dataHttpd := templates.HttpdTemplateData{
		Volumes:       utils.PROXY_HTTPD_VOLUMES,
		HttpProxyFile: httpProxyConfig,
		Image:         httpdImage,
	}
	if err := generateSystemdFile(dataHttpd, "httpd"); err != nil {
		return fmt.Errorf("cannot generated systemd file: %s", err)
	}

	// Salt broker
	dataSaltBroker := templates.SaltBrokerTemplateData{
		HttpProxyFile: httpProxyConfig,
		Image:         saltBrokerImage,
	}
	if err := generateSystemdFile(dataSaltBroker, "salt-broker"); err != nil {
		return fmt.Errorf("cannot generated systemd file: %s", err)
	}

	// Squid
	dataSquid := templates.SquidTemplateData{
		Volumes:       utils.PROXY_SQUID_VOLUMES,
		HttpProxyFile: httpProxyConfig,
		Image:         squidImage,
	}
	if err := generateSystemdFile(dataSquid, "squid"); err != nil {
		return fmt.Errorf("cannot generated systemd file: %s", err)
	}

	// SSH
	dataSSH := templates.SSHTemplateData{
		HttpProxyFile: httpProxyConfig,
		Image:         sshImage,
	}
	if err := generateSystemdFile(dataSSH, "ssh"); err != nil {
		return fmt.Errorf("cannot generated systemd file: %s", err)
	}

	// Tftpd
	dataTftpd := templates.TFTPDTemplateData{
		Volumes:       utils.PROXY_TFTPD_VOLUMES,
		HttpProxyFile: httpProxyConfig,
		Image:         tftpdImage,
	}
	if err := generateSystemdFile(dataTftpd, "tftpd"); err != nil {
		return fmt.Errorf("cannot generated systemd file: %s", err)
	}

	return podman.ReloadDaemon(false)
}

func generateSystemdFile(template utils.Template, service string) error {
	name := fmt.Sprintf("uyuni-proxy-%s.service", service)

	const systemdPath = "/etc/systemd/system"
	path := path.Join(systemdPath, name)
	if err := utils.WriteTemplateToFile(template, path, 0555, true); err != nil {
		return fmt.Errorf("failed to generate %s", path)
	}
	return nil
}

func getHttpProxyConfig() string {
	const httpProxyConfigPath = "/etc/sysconfig/proxy"

	// Only SUSE distros seem to have such a file for HTTP proxy settings
	if utils.FileExists(httpProxyConfigPath) {
		return httpProxyConfigPath
	}
	return ""
}
