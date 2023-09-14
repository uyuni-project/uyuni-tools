package podman

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

const UYUNI_NETWORK = "uyuni"
const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

func GetCommonParams(containerName string) []string {
	return strings.Split(fmt.Sprintf(commonArgs, containerName), " ")
}

func GetExposedPorts() []string {

	ports := []string{"443", "80"}
	for _, portMap := range utils.TCP_PORTS {
		ports = append(ports, strconv.Itoa(portMap.Port))
	}
	for _, portMap := range utils.UDP_PORTS {
		ports = append(ports, strconv.Itoa(portMap.Port))
	}

	return ports
}

const ServicePath = "/etc/systemd/system/uyuni-server.service"

func GenerateSystemdService(tz string, image string, podmanArgs []string) {

	setupNetwork()

	data := templates.PodmanServiceTemplateData{
		Volumes:    utils.VOLUMES,
		NamePrefix: "uyuni",
		Args:       commonArgs + " " + strings.Join(podmanArgs, " "),
		Ports:      GetExposedPorts(),
		Timezone:   tz,
		Image:      image,
		Network:    UYUNI_NETWORK,
	}
	if err := utils.WriteTemplateToFile(data, ServicePath, 0555, false); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate systemd service unit file")
	}
	utils.RunRawCmd("systemctl", []string{"daemon-reload"}, true)
}

func setupNetwork() {
	log.Info().Msgf("Setting up %s network", UYUNI_NETWORK)

	ipv6Enabled := isIpv6Enabled()

	// Check if the uyuni network exists and is IPv6 enabled
	hasIpv6, err := exec.Command("podman", "network", "inspect", "--format", "{{.IPv6Enabled}}", UYUNI_NETWORK).Output()
	if err == nil {
		if string(hasIpv6) != "true" && ipv6Enabled {
			log.Info().Msgf("%s network doesn't have IPv6, deleting existing network to enable IPv6 on it", UYUNI_NETWORK)
			message := fmt.Sprintf("Failed to remove %s podman network", UYUNI_NETWORK)
			err := utils.RunRawCmd("podman", []string{"network", "rm", UYUNI_NETWORK,
				"--log-level", log.Logger.GetLevel().String()}, true)
			if err != nil {
				log.Fatal().Err(err).Msg(message)
			}
		} else {
			log.Info().Msgf("Reusing existing %s network", UYUNI_NETWORK)
			return
		}
	}

	message := fmt.Sprintf("Failed to create %s network with IPv6 enabled", UYUNI_NETWORK)

	args := []string{"network", "create"}
	if ipv6Enabled {
		// An IPv6 network on a host where IPv6 is disabled doesn't work: don't try it.
		// Check if the networkd backend is netavark
		out, err := exec.Command("podman", "info", "--format", "{{.Host.NetworkBackend}}").Output()
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
	err = utils.RunRawCmd("podman", args, true)
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
	out, err := os.ReadFile(file)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to read file %s", file)
	}
	return string(out[:]) != "0"
}

func EnablePodmanSocket() {
	err := utils.RunRawCmd("systemctl", []string{"enable", "--now", "podman.socket"}, true)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to enable podman.socket unit")
	}
}
