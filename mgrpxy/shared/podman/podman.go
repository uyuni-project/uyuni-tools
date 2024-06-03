// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
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
	httpProxyConfig := getHttpProxyConfig()

	ports := []types.PortMap{}
	ports = append(ports, shared_utils.PROXY_TCP_PORTS...)
	ports = append(ports, shared_utils.PROXY_PODMAN_PORTS...)
	ports = append(ports, shared_utils.UDP_PORTS...)

	// Pod
	dataPod := templates.PodTemplateData{
		Ports:         ports,
		HttpProxyFile: httpProxyConfig,
		Network:       podman.UyuniNetwork,
		IPV6Enabled:   ipv6Enabled,
	}
	podEnv := fmt.Sprintf(`Environment="PODMAN_EXTRA_ARGS=%s"`, strings.Join(flags.Podman.Args, " "))
	if err := generateSystemdFile(dataPod, "pod", "", podEnv); err != nil {
		return err
	}

	// Httpd
	{
		dataHttpd := templates.HttpdTemplateData{
			Volumes:       shared_utils.PROXY_HTTPD_VOLUMES,
			HttpProxyFile: httpProxyConfig,
		}
		additionHttpdTuningSettings := ""
		if flags.ProxyImageFlags.Tuning.Httpd != "" {
			absPath, err := filepath.Abs(flags.ProxyImageFlags.Tuning.Httpd)
			if err != nil {
				return err
			}
			additionHttpdTuningSettings = fmt.Sprintf(`Environment=HTTPD_EXTRA_CONF=-v%s:/etc/apache2/conf.d/apache_tuning.conf:ro`, absPath)
		}
		if err := generateSystemdFile(dataHttpd, "httpd", httpdImage, additionHttpdTuningSettings); err != nil {
			return err
		}
	}
	// Salt broker
	{
		dataSaltBroker := templates.SaltBrokerTemplateData{
			HttpProxyFile: httpProxyConfig,
		}
		if err := generateSystemdFile(dataSaltBroker, "salt-broker", saltBrokerImage, ""); err != nil {
			return err
		}
	}
	// Squid
	{
		dataSquid := templates.SquidTemplateData{
			Volumes:       shared_utils.PROXY_SQUID_VOLUMES,
			HttpProxyFile: httpProxyConfig,
		}
		additionSquidTuningSettings := ""
		if flags.ProxyImageFlags.Tuning.Squid != "" {
			absPath, err := filepath.Abs(flags.ProxyImageFlags.Tuning.Squid)
			if err != nil {
				return err
			}
			additionSquidTuningSettings = fmt.Sprintf(`Environment=SQUID_EXTRA_CONF=-v%s:/etc/squid/conf.d/squid_tuning.conf:ro`, absPath)
		}
		if err := generateSystemdFile(dataSquid, "squid", squidImage, additionSquidTuningSettings); err != nil {
			return err
		}
	}
	// SSH
	{
		dataSSH := templates.SSHTemplateData{
			HttpProxyFile: httpProxyConfig,
		}
		if err := generateSystemdFile(dataSSH, "ssh", sshImage, ""); err != nil {
			return err
		}
	}
	// Tftpd
	{
		dataTftpd := templates.TFTPDTemplateData{
			Volumes:       shared_utils.PROXY_TFTPD_VOLUMES,
			HttpProxyFile: httpProxyConfig,
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
		if err := podman.GenerateSystemdConfFile("uyuni-proxy-"+service, "custom.conf", config, false); err != nil {
			return shared_utils.Errorf(err, L("cannot generate systemd conf user configuration file"))
		}
	}
	return nil
}

func getHttpProxyConfig() string {
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
	log.Info().Msgf(L("Setting up proxy with configuration %s"), configPath)
	const proxyConfigDir = "/etc/uyuni/proxy"
	if err := os.MkdirAll(proxyConfigDir, 755); err != nil {
		return err
	}

	if err := shared_utils.ExtractTarGz(configPath, proxyConfigDir); err != nil {
		return err
	}

	proxyConfigDirInfo, err := os.Stat(proxyConfigDir)
	if err != nil {
		return err
	}

	dirMode := proxyConfigDirInfo.Mode()

	if !(dirMode&0005 != 0 && dirMode&0050 != 0 && dirMode&0500 != 0) {
		return fmt.Errorf(L("/etc/uyuni/proxy directory has no read and write permissions for all users. Check your umask settings."))
	}

	if err := shared_utils.ExtractTarGz(configPath, proxyConfigDir); err != nil {
		return err
	}

	proxyConfigInfo, err := os.Stat(path.Join(proxyConfigDir, "config.yaml"))
	if err != nil {
		return err
	}

	mode := proxyConfigInfo.Mode()

	if !(mode&0004 != 0 && mode&0040 != 0 && mode&0400 != 0) {
		return fmt.Errorf(L("/etc/uyuni/proxy/config.yaml has no read permissions for all users. Check your umask settings."))
	}

	return nil
}

// Upgrade will upgrade the proxy podman deploy.
func Upgrade(
	systemd podman.Systemd, globalFlags *types.GlobalFlags, flags *PodmanProxyFlags,
	cmd *cobra.Command, args []string,
) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf(L("install podman before running this command"))
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
	} else {
		return systemd.EnableService(podman.ProxyService)
	}
}
