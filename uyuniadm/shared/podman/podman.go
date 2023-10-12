package podman

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

const UYUNI_NETWORK = "uyuni"
const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

func GetCommonParams() []string {
	return strings.Split(commonArgs, " ")
}

func GetExposedPorts(debug bool) []utils.PortMap {

	ports := []utils.PortMap{
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

const ServicePath = "/etc/systemd/system/uyuni-server.service"

func GenerateSystemdService(tz string, image string, debug bool, podmanArgs []string) {

	setupNetwork()

	log.Info().Msg("Enabling system service")
	data := templates.PodmanServiceTemplateData{
		Volumes:    utils.VOLUMES,
		NamePrefix: "uyuni",
		Args:       commonArgs + " " + strings.Join(podmanArgs, " "),
		Ports:      GetExposedPorts(debug),
		Timezone:   tz,
		Image:      image,
		Network:    UYUNI_NETWORK,
	}
	if err := utils.WriteTemplateToFile(data, ServicePath, 0555, false); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate systemd service unit file")
	}

	utils.RunCmd("systemctl", "daemon-reload")
}

func setupNetwork() {
	log.Info().Msgf("Setting up %s network", UYUNI_NETWORK)

	ipv6Enabled := isIpv6Enabled()

	testNetworkCmd := exec.Command("podman", "network", "exists", UYUNI_NETWORK)
	testNetworkCmd.Run()
	// check if network exists before trying to get the IPV6 information
	if testNetworkCmd.ProcessState.ExitCode() == 0 {
		// Check if the uyuni network exists and is IPv6 enabled
		hasIpv6, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "network", "inspect", "--format", "{{.IPv6Enabled}}", UYUNI_NETWORK)
		if err == nil {
			if string(hasIpv6) != "true" && ipv6Enabled {
				log.Info().Msgf("%s network doesn't have IPv6, deleting existing network to enable IPv6 on it", UYUNI_NETWORK)
				message := fmt.Sprintf("Failed to remove %s podman network", UYUNI_NETWORK)
				err := utils.RunCmd("podman", "network", "rm", UYUNI_NETWORK,
					"--log-level", log.Logger.GetLevel().String())
				if err != nil {
					log.Fatal().Err(err).Msg(message)
				}
			} else {
				log.Info().Msgf("Reusing existing %s network", UYUNI_NETWORK)
				return
			}
		}
	}

	message := fmt.Sprintf("Failed to create %s network with IPv6 enabled", UYUNI_NETWORK)

	args := []string{"network", "create"}
	if ipv6Enabled {
		// An IPv6 network on a host where IPv6 is disabled doesn't work: don't try it.
		// Check if the networkd backend is netavark
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "info", "--format", "{{.Host.NetworkBackend}}")
		backend := strings.Trim(string(out), "\n")
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to find podman's network backend")
		} else if backend != "netavark" {
			log.Info().Msgf("Podman's network backend (%s) is not netavark, skipping IPv6 enabling on %s network", backend, UYUNI_NETWORK)
		} else {
			args = append(args, "--ipv6")
		}
	}
	args = append(args, UYUNI_NETWORK)
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
		if getFileBoolean(file) {
			return false
		}
	}
	return true
}

func getFileBoolean(file string) bool {
	return string(utils.ReadFile(file)) != "0"
}

func EnablePodmanSocket() {
	err := utils.RunCmd("systemctl", "enable", "--now", "podman.socket")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to enable podman.socket unit")
	}
}

func UpdateSslCertificate(cnx *utils.Connection, chain *ssl.CaChain, serverPair *ssl.SslPair) {
	ssl.CheckPaths(chain, serverPair)

	// Copy the CAs, certificate and key to the container
	const certDir = "/tmp/uyuni-tools"
	if err := utils.RunCmd("podman", "exec", utils.PODMAN_CONTAINER, "mkdir", "-p", certDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to create temporary folder on container to copy certificates to")
	}

	rootCaPath := path.Join(certDir, "root-ca.crt")
	serverCrtPath := path.Join(certDir, "server.crt")
	serverKeyPath := path.Join(certDir, "server.key")

	log.Debug().Msgf("Intermediate CA flags: %v", chain.Intermediate)

	args := []string{
		"exec",
		utils.PODMAN_CONTAINER,
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
	if err := utils.RunCmd("podman", "exec", utils.PODMAN_CONTAINER, "rm", "-rf", certDir); err != nil {
		log.Error().Err(err).Msg("Failed to remove copied certificate files in the container")
	}

	const sslbuildPath = "/root/ssl-build"
	if utils.TestExistenceInPod(cnx, sslbuildPath) {
		if err := utils.RunCmd("podman", "exec", utils.PODMAN_CONTAINER, "rm", "-rf", sslbuildPath); err != nil {
			log.Error().Err(err).Msg("Failed to remove now useless ssl-build folder in the container")
		}
	}

	// The services need to be restarted
	log.Info().Msg("Restarting services after updating the certificate")
	utils.RunCmdStdMapping("podman", "exec", utils.PODMAN_CONTAINER, "spacewalk-service", "restart")
}
