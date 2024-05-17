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
	Podman                podman.PodmanFlags
}

// GenerateSystemdService generates all the systemd files required by proxy.
func GenerateSystemdService(httpdImage string, saltBrokerImage string, squidImage string, sshImage string,
	tftpdImage string, flags *PodmanProxyFlags) error {
	if err := podman.SetupNetwork(true); err != nil {
		return shared_utils.Errorf(err, L("cannot setup network"))
	}

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
	}
	podEnv := fmt.Sprintf(`Environment="PODMAN_EXTRA_ARGS=%s"`, strings.Join(podmanArgs, " "))
	if err := generateSystemdFile(dataPod, "pod", "", podEnv); err != nil {
		return err
	}

	httpdVolumes := shared_utils.PROXY_HTTPD_VOLUMES
	log.Debug().Msgf("Tuning HTTPD value: %s", flags.ProxyImageFlags.TuningHttpd)
	if flags.ProxyImageFlags.TuningHttpd != "" {
		absPath, err := filepath.Abs(flags.ProxyImageFlags.TuningHttpd)
		if err != nil {
			return err
		}
		if !shared_utils.FileExists(absPath) {
			return fmt.Errorf(L("%s does not exists"), absPath)
		}
		httpdVolumes[absPath] = "/etc/apache2/conf.d/apache_tuning.conf"
	}
	// Httpd
	dataHttpd := templates.HttpdTemplateData{
		Volumes:       httpdVolumes,
		HttpProxyFile: httpProxyConfig,
	}
	if err := generateSystemdFile(dataHttpd, "httpd", httpdImage, ""); err != nil {
		return err
	}
	if flags.ProxyImageFlags.TuningHttpd != "" {
		absPath, err := filepath.Abs(flags.ProxyImageFlags.TuningHttpd)
		if err != nil {
			return err
		}
		additionPodmanSettings := fmt.Sprintf(`Environment=HTTPD_ADDITIONAL_SETTINGS=-v%s:/etc/apache2/conf.d/apache_tuning.conf:ro`, absPath)
		if err := podman.GenerateSystemdConfFile("uyuni-proxy-httpd", "Service", additionPodmanSettings); err != nil {
			return shared_utils.Errorf(err, L("cannot generate systemd conf file"))
		}
	}

	// Salt broker
	dataSaltBroker := templates.SaltBrokerTemplateData{
		HttpProxyFile: httpProxyConfig,
	}
	if err := generateSystemdFile(dataSaltBroker, "salt-broker", saltBrokerImage, ""); err != nil {
		return err
	}

	squidVolumes := shared_utils.PROXY_SQUID_VOLUMES

	log.Debug().Msgf("Tuning Suid value: %s", flags.ProxyImageFlags.TuningSquid)
	if flags.ProxyImageFlags.TuningSquid != "" {
		absPath, err := filepath.Abs(flags.ProxyImageFlags.TuningSquid)
		if err != nil {
			return err
		}
		additionPodmanSettings := fmt.Sprintf(`Environment=SQUID_ADDITIONAL_SETTINGS=-v%s:/etc/squid/conf.d/squid_tuning.conf:ro`, absPath)
		if err := podman.GenerateSystemdConfFile("uyuni-proxy-squid", "Service", additionPodmanSettings); err != nil {
			return shared_utils.Errorf(err, L("cannot generate systemd conf file"))
		}
	}

	// Squid
	dataSquid := templates.SquidTemplateData{
		Volumes:       squidVolumes,
		HttpProxyFile: httpProxyConfig,
	}
	if err := generateSystemdFile(dataSquid, "squid", squidImage, ""); err != nil {
		return err
	}

	// SSH
	dataSSH := templates.SSHTemplateData{
		HttpProxyFile: httpProxyConfig,
	}
	if err := generateSystemdFile(dataSSH, "ssh", sshImage, ""); err != nil {
		return err
	}

	// Tftpd
	dataTftpd := templates.TFTPDTemplateData{
		Volumes:       shared_utils.PROXY_TFTPD_VOLUMES,
		HttpProxyFile: httpProxyConfig,
	}
	if err := generateSystemdFile(dataTftpd, "tftpd", tftpdImage, ""); err != nil {
		return err
	}

	return podman.ReloadDaemon(false)
}

func generateSystemdFile(template shared_utils.Template, service string, image string, config string) error {
	name := fmt.Sprintf("uyuni-proxy-%s.service", service)

	const systemdPath = "/etc/systemd/system"
	path := path.Join(systemdPath, name)
	if err := shared_utils.WriteTemplateToFile(template, path, 0644, true); err != nil {
		return shared_utils.Errorf(err, L("failed to generate systemd file '%s'"), path)
	}

	if image != "" || config != "" {
		configBody := fmt.Sprintf(`%s
Environment=UYUNI_IMAGE=%s`, config, image)
		if err := podman.GenerateSystemdConfFile("uyuni-proxy-"+service, "Service", configBody); err != nil {
			return shared_utils.Errorf(err, L("cannot generate systemd conf file"))
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
func GetContainerImage(flags *utils.ProxyImageFlags, name string) (string, error) {
	image := flags.GetContainerImage(name)
	inspectedHostValues, err := shared_utils.InspectHost(true)
	if err != nil {
		return "", shared_utils.Errorf(err, L("cannot inspect host values"))
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := podman.PrepareImage(image, flags.PullPolicy, pullArgs...)
	if err != nil {
		return "", err
	}

	return preparedImage, nil
}

// UnpackConfig uncompress the config.tar.gz containing proxy configuration.
func UnpackConfig(configPath string) error {
	log.Info().Msgf(L("Setting up proxy with configuration %s"), configPath)
	const proxyConfigDir = "/etc/uyuni/proxy"
	if err := os.MkdirAll(proxyConfigDir, 0755); err != nil {
		return err
	}

	if err := shared_utils.ExtractTarGz(configPath, proxyConfigDir); err != nil {
		return err
	}
	return nil
}

// Upgrade will upgrade the proxy podman deploy.
func Upgrade(globalFlags *types.GlobalFlags, flags *PodmanProxyFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf(L("install podman before running this command"))
	}

	httpdImage, err := getContainerImage(&flags.ProxyImageFlags, "httpd")
	if err != nil {
		log.Info().Msgf(L("cannot find httpd image: it will no be upgraded"))
	}
	saltBrokerImage, err := getContainerImage(&flags.ProxyImageFlags, "salt-broker")
	if err != nil {
		log.Info().Msgf(L("cannot find salt-broker image: it will no be upgraded"))
	}
	squidImage, err := getContainerImage(&flags.ProxyImageFlags, "squid")
	if err != nil {
		log.Info().Msgf(L("cannot find squid image: it will no be upgraded"))
	}
	sshImage, err := getContainerImage(&flags.ProxyImageFlags, "ssh")
	if err != nil {
		log.Info().Msgf(L("cannot find ssh image: it will no be upgraded"))
	}
	tftpdImage, err := getContainerImage(&flags.ProxyImageFlags, "tftpd")
	if err != nil {
		log.Info().Msgf(L("cannot find tftpd image: it will no be upgraded"))
	}

	// Setup the systemd service configuration options
	if err := GenerateSystemdService(httpdImage, saltBrokerImage, squidImage, sshImage, tftpdImage, flags); err != nil {
		return err
	}

	return startPod()
}

func getContainerImage(flags *utils.ProxyImageFlags, name string) (string, error) {
	image := flags.GetContainerImage(name)
	inspectedHostValues, err := shared_utils.InspectHost(true)
	if err != nil {
		return "", shared_utils.Errorf(err, L("cannot inspect host values"))
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := podman.PrepareImage(image, flags.PullPolicy, pullArgs...)
	if err != nil {
		return "", err
	}

	return preparedImage, nil
}

// Start the proxy services.
func startPod() error {
	ret := podman.IsServiceRunning(podman.ProxyService)
	if ret {
		return podman.RestartService(podman.ProxyService)
	} else {
		return podman.EnableService(podman.ProxyService)
	}
}
