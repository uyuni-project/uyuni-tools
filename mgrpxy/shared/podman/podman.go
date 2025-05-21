// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/templates"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

// PodmanProxyFlags are the flags used by podman proxy install and upgrade command.
type PodmanProxyFlags struct {
	utils.ProxyImageFlags `mapstructure:",squash"`
	SCC                   types.SCCCredentials
	Podman                podman.PodmanFlags
}

// GenerateSystemdService generates all the systemd files required by proxy.
func GenerateSystemdService(
	systemd podman.Systemd,
	httpdImage string,
	saltBrokerImage string,
	squidImage string,
	sshImage string,
	tftpdImage string,
	flags *PodmanProxyFlags,
) error {
	err := podman.SetupNetwork(true)
	if err != nil {
		return shared_utils.Errorf(err, L("cannot setup network"))
	}

	ipv6Enabled := podman.HasIpv6Enabled(podman.UyuniNetwork)

	log.Info().Msg(L("Generating systemd services"))
	httpProxyConfig := getHTTPProxyConfig()

	ports := []types.PortMap{}
	ports = append(ports, shared_utils.ProxyTCPPorts...)
	ports = append(ports, shared_utils.ProxyPodmanPorts...)
	ports = append(ports, shared_utils.TftpPorts...)

	// Pod
	dataPod := templates.PodTemplateData{
		Ports:         ports,
		HTTPProxyFile: httpProxyConfig,
		Network:       podman.UyuniNetwork,
		IPV6Enabled:   ipv6Enabled,
	}
	podEnv := fmt.Sprintf(`Environment="PODMAN_EXTRA_ARGS=%s"`, strings.Join(flags.Podman.Args, " "))
	if err := generateSystemdFile(dataPod, "pod", "", podEnv); err != nil {
		return err
	}

	// Httpd
	volumeOptions := ""

	{
		dataHttpd := templates.HttpdTemplateData{
			Volumes:       shared_utils.ProxyHttpdVolumes,
			HTTPProxyFile: httpProxyConfig,
		}
		additionHttpdTuningSettings := ""
		if flags.ProxyImageFlags.Tuning.Httpd != "" {
			absPath, err := filepath.Abs(flags.ProxyImageFlags.Tuning.Httpd)
			if err != nil {
				return err
			}
			additionHttpdTuningSettings = fmt.Sprintf(
				`Environment=HTTPD_EXTRA_CONF=-v%s:/etc/apache2/conf.d/apache_tuning.conf:ro%s`,
				absPath, volumeOptions,
			)
		}
		if err := generateSystemdFile(dataHttpd, "httpd", httpdImage, additionHttpdTuningSettings); err != nil {
			return err
		}
	}
	// Salt broker
	{
		dataSaltBroker := templates.SaltBrokerTemplateData{
			HTTPProxyFile: httpProxyConfig,
		}
		if err := generateSystemdFile(dataSaltBroker, "salt-broker", saltBrokerImage, ""); err != nil {
			return err
		}
	}
	// Squid
	{
		dataSquid := templates.SquidTemplateData{
			Volumes:       shared_utils.ProxySquidVolumes,
			HTTPProxyFile: httpProxyConfig,
		}
		additionSquidTuningSettings := ""
		if flags.ProxyImageFlags.Tuning.Squid != "" {
			absPath, err := filepath.Abs(flags.ProxyImageFlags.Tuning.Squid)
			if err != nil {
				return err
			}
			additionSquidTuningSettings = fmt.Sprintf(
				`Environment=SQUID_EXTRA_CONF=-v%s:/etc/squid/conf.d/squid_tuning.conf:ro%s`,
				absPath, volumeOptions,
			)
		}
		if err := generateSystemdFile(dataSquid, "squid", squidImage, additionSquidTuningSettings); err != nil {
			return err
		}
	}
	// SSH
	{
		dataSSH := templates.SSHTemplateData{
			HTTPProxyFile: httpProxyConfig,
		}
		if err := generateSystemdFile(dataSSH, "ssh", sshImage, ""); err != nil {
			return err
		}
	}
	// Tftpd
	{
		dataTftpd := templates.TFTPDTemplateData{
			Volumes:       shared_utils.ProxyTftpdVolumes,
			HTTPProxyFile: httpProxyConfig,
		}
		if err := generateSystemdFile(dataTftpd, "tftpd", tftpdImage, ""); err != nil {
			return err
		}
	}
	return systemd.ReloadDaemon(false)
}

func generateSystemdFile(template shared_utils.Template, service string, image string, config string) error {
	name := fmt.Sprintf("uyuni-proxy-%s.service", service)

	const systemdPath = "/etc/systemd/system"
	path := path.Join(systemdPath, name)
	if err := shared_utils.WriteTemplateToFile(template, path, 0644, true); err != nil {
		return shared_utils.Errorf(err, L("failed to generate systemd file '%s'"), path)
	}

	if image != "" {
		configBody := fmt.Sprintf("Environment=UYUNI_IMAGE=%s", image)
		if err := podman.GenerateSystemdConfFile("uyuni-proxy-"+service, "generated.conf", configBody, true); err != nil {
			return shared_utils.Errorf(err, L("cannot generate systemd conf file"))
		}
	}

	if config != "" {
		if err := podman.GenerateSystemdConfFile("uyuni-proxy-"+service, podman.CustomConf, config, false); err != nil {
			return shared_utils.Errorf(err, L("cannot generate systemd conf user configuration file"))
		}
	}
	return nil
}

func getHTTPProxyConfig() string {
	const httpProxyConfigPath = "/etc/sysconfig/proxy"

	// Only SUSE distros seem to have such a file for HTTP proxy settings
	if shared_utils.FileExists(httpProxyConfigPath) {
		return httpProxyConfigPath
	}
	return ""
}

// GetContainerImage returns a proxy image URL.
func GetContainerImage(authFile string, flags *utils.ProxyImageFlags, name string) (string, error) {
	image := flags.GetContainerImage(name)

	preparedImage, err := podman.PrepareImage(authFile, image, flags.PullPolicy, true)
	if err != nil {
		return "", err
	}

	return preparedImage, nil
}

// UnpackConfig uncompress the config.tar.gz containing proxy configuration.
func UnpackConfig(configPath string) error {
	const proxyConfigDir = "/etc/uyuni/proxy"

	// Create dir if it doesn't exist & check perms
	if err := os.MkdirAll(proxyConfigDir, 0755); err != nil {
		return err
	}

	if err := checkPermissions(proxyConfigDir, 0005|0050|0500); err != nil {
		return err
	}

	// Extract the tarball, if provided
	if configPath != "" {
		log.Info().Msgf(L("Setting up proxy with configuration %s"), configPath)
		if err := shared_utils.ExtractTarGz(configPath, proxyConfigDir); err != nil {
			return shared_utils.Errorf(err, L("failed to extract proxy config from %s file"), configPath)
		}
	} else {
		log.Info().Msg(L("No tarball provided. Will check existing configuration files."))
	}

	return validateInstallYamlFiles(proxyConfigDir)
}

// checkPermissions checks if a directory or file has a required permissions.
func checkPermissions(path string, requiredMode os.FileMode) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Mode()&requiredMode != requiredMode {
		if info.IsDir() {
			return fmt.Errorf(L("%s directory has no required permissions. Check your umask settings"), path)
		}
		return fmt.Errorf(L("%s file has no required permissions. Check your umask settings"), path)
	}
	return nil
}

// validateYamlFiles validates if the required configuration files.
func validateInstallYamlFiles(dir string) error {
	yamlFiles := []string{"httpd.yaml", "ssh.yaml", "config.yaml"}

	for _, file := range yamlFiles {
		filePath := path.Join(dir, file)
		_, err := os.Stat(filePath)
		if err != nil {
			return fmt.Errorf(L("missing required configuration file: %s"), filePath)
		}
		if file == "config.yaml" {
			if err := checkPermissions(filePath, 0004|0040|0400); err != nil {
				return err
			}
		}
	}
	return nil
}

// Upgrade will upgrade the proxy podman deploy.
func Upgrade(
	systemd podman.Systemd, _ *types.GlobalFlags, flags *PodmanProxyFlags,
	_ *cobra.Command, _ []string,
) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}
	if err := systemd.StopService(podman.ProxyService); err != nil {
		return err
	}

	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman.PodmanLogin(hostData, flags.SCC)
	if err != nil {
		return shared_utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	httpdImage, err := GetContainerImage(authFile, &flags.ProxyImageFlags, "httpd")
	if err != nil {
		log.Warn().Msgf(L("cannot find httpd image: it will no be upgraded"))
	}
	saltBrokerImage, err := GetContainerImage(authFile, &flags.ProxyImageFlags, "salt-broker")
	if err != nil {
		log.Warn().Msgf(L("cannot find salt-broker image: it will no be upgraded"))
	}
	squidImage, err := GetContainerImage(authFile, &flags.ProxyImageFlags, "squid")
	if err != nil {
		log.Warn().Msgf(L("cannot find squid image: it will no be upgraded"))
	}
	sshImage, err := GetContainerImage(authFile, &flags.ProxyImageFlags, "ssh")
	if err != nil {
		log.Warn().Msgf(L("cannot find ssh image: it will no be upgraded"))
	}
	tftpdImage, err := GetContainerImage(authFile, &flags.ProxyImageFlags, "tftpd")
	if err != nil {
		log.Warn().Msgf(L("cannot find tftpd image: it will no be upgraded"))
	}

	// Setup the systemd service configuration options
	err = GenerateSystemdService(systemd, httpdImage, saltBrokerImage, squidImage, sshImage, tftpdImage, flags)
	if err != nil {
		return err
	}

	return startPod(systemd)
}

// Start the proxy services.
func startPod(systemd podman.Systemd) error {
	ret := systemd.IsServiceRunning(podman.ProxyService)
	if ret {
		return systemd.RestartService(podman.ProxyService)
	}
	return systemd.EnableService(podman.ProxyService)
}
