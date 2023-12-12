// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

func GetCommonParams() []string {
	return strings.Split(commonArgs, " ")
}

func GetExposedPorts(debug bool) []types.PortMap {

	ports := []types.PortMap{
		utils.NewPortMap("https", 443, 443),
		utils.NewPortMap("http", 80, 80),
	}
	ports = append(ports, utils.TCP_PORTS...)
	ports = append(ports, utils.UDP_PORTS...)

	if debug {
		ports = append(ports, utils.DEBUG_PORTS...)
	}

	return ports
}

func GenerateSystemdService(tz string, image string, debug bool, podmanArgs []string) {

	podman.SetupNetwork()

	log.Info().Msg("Enabling system service")
	data := templates.PodmanServiceTemplateData{
		Volumes:    utils.VOLUMES,
		NamePrefix: "uyuni",
		Args:       commonArgs + " " + strings.Join(podmanArgs, " "),
		Ports:      GetExposedPorts(debug),
		Timezone:   tz,
		Image:      image,
		Network:    podman.UyuniNetwork,
	}
	if err := utils.WriteTemplateToFile(data, podman.GetServicePath("uyuni-server"), 0555, false); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate systemd service unit file")
	}

	utils.RunCmd("systemctl", "daemon-reload")
}

func UpdateSslCertificate(cnx *utils.Connection, chain *ssl.CaChain, serverPair *ssl.SslPair) {
	ssl.CheckPaths(chain, serverPair)

	// Copy the CAs, certificate and key to the container
	const certDir = "/tmp/uyuni-tools"
	if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "mkdir", "-p", certDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to create temporary folder on container to copy certificates to")
	}

	rootCaPath := path.Join(certDir, "root-ca.crt")
	serverCrtPath := path.Join(certDir, "server.crt")
	serverKeyPath := path.Join(certDir, "server.key")

	log.Debug().Msgf("Intermediate CA flags: %v", chain.Intermediate)

	args := []string{
		"exec",
		podman.ServerContainerName,
		"mgr-ssl-cert-setup",
		"-vvv",
		"--root-ca-file", rootCaPath,
		"--server-cert-file", serverCrtPath,
		"--server-key-file", serverKeyPath,
	}

	utils.Copy(cnx, chain.Root, "server:"+rootCaPath, "root", "root")
	utils.Copy(cnx, serverPair.Cert, "server:"+serverCrtPath, "root", "root")
	utils.Copy(cnx, serverPair.Key, "server:"+serverKeyPath, "root", "root")

	for i, ca := range chain.Intermediate {
		caFilename := fmt.Sprintf("ca-%d.crt", i)
		caPath := path.Join(certDir, caFilename)
		args = append(args, "--intermediate-ca-file", caPath)
		utils.Copy(cnx, ca, "server:"+caPath, "root", "root")
	}

	// Check and install then using mgr-ssl-cert-setup
	if _, err := utils.RunCmdOutput(zerolog.InfoLevel, "podman", args...); err != nil {
		log.Fatal().Err(err).Msg("Failed to update SSL certificate")
	}

	// Clean the copied files and the now useless ssl-build
	if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "rm", "-rf", certDir); err != nil {
		log.Error().Err(err).Msg("Failed to remove copied certificate files in the container")
	}

	const sslbuildPath = "/root/ssl-build"
	if utils.TestExistenceInPod(cnx, sslbuildPath) {
		if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "rm", "-rf", sslbuildPath); err != nil {
			log.Error().Err(err).Msg("Failed to remove now useless ssl-build folder in the container")
		}
	}

	// The services need to be restarted
	log.Info().Msg("Restarting services after updating the certificate")
	utils.RunCmdStdMapping("podman", "exec", podman.ServerContainerName, "spacewalk-service", "restart")
}
