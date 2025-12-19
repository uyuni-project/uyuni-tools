// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/templates"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

const (
	SystemIDEvent         = "suse/systemid/generate"
	SystemIDEventResponse = "suse/systemid/generated"
	SystemIDSecret        = "uyuni-proxy-systemid"
	defaultApacheConf     = "/etc/uyuni/proxy/apache.conf"
	defaultSquidConf      = "/etc/uyuni/proxy/squid.conf"
	defaultSSHConf        = "/etc/uyuni/proxy/ssh.conf"
	ServiceHTTPd          = "uyuni-proxy-httpd"
	ServiceSSH            = "uyuni-proxy-ssh"
	ServiceSquid          = "uyuni-proxy-squid"
	ServiceSaltBroker     = "uyuni-proxy-salt-broker"
	ServiceTFTFd          = "uyuni-proxy-tftpd"
)

var contextRunner = shared_utils.NewRunnerWithContext
var newRunner = shared_utils.NewRunner

// PodmanProxyFlags are the flags used by podman proxy install and upgrade command.
type PodmanProxyFlags struct {
	utils.ProxyImageFlags `mapstructure:",squash"`
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
		if podman.HasSecret(SystemIDSecret) {
			dataHttpd.SystemIDSecret = SystemIDSecret
		}

		additionHttpdTuningSettings := ""
		additionHTTPConfPath, err := getPathOrDefault(flags.ProxyImageFlags.Tuning.Httpd, defaultApacheConf)
		if err != nil {
			return err
		}

		if additionHTTPConfPath != "" {
			additionHttpdTuningSettings = fmt.Sprintf(
				`Environment=HTTPD_EXTRA_CONF=-v%s:/etc/apache2/conf.d/apache_tuning.conf:ro%s`,
				additionHTTPConfPath, volumeOptions,
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
		additionSquidConfPath, err := getPathOrDefault(flags.ProxyImageFlags.Tuning.Squid, defaultSquidConf)
		if err != nil {
			return err
		}
		if additionSquidConfPath != "" {
			additionSquidTuningSettings = fmt.Sprintf(
				`Environment=SQUID_EXTRA_CONF=-v%s:/etc/squid/conf.d/squid_tuning.conf:ro%s`,
				additionSquidConfPath, volumeOptions,
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
		additionSSHTuningSettings := ""
		additionSSHConfPath, err := getPathOrDefault(flags.ProxyImageFlags.Tuning.SSH, defaultSSHConf)
		if err != nil {
			return err
		}
		if additionSSHConfPath != "" {
			additionSSHTuningSettings = fmt.Sprintf(
				`Environment=SSH_EXTRA_CONF=-v%s:/etc/ssh/sshd_config.d/99-tuning.conf:ro%s`,
				additionSSHConfPath, volumeOptions,
			)
		}
		if err := generateSystemdFile(dataSSH, "ssh", sshImage, additionSSHTuningSettings); err != nil {
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

func getPathOrDefault(flag string, defaultPath string) (string, error) {
	result := ""
	if flag != "" {
		var err error
		result, err = filepath.Abs(flag)
		if err != nil {
			return "", err
		}
	} else if shared_utils.FileExists(defaultPath) {
		result = defaultPath
	}

	return result, nil
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

	// Check if we are a salt minion registered to SMLM and if so, try to get up to date systemid
	if hostData.HasSaltMinion {
		// If we previously created systemid secret, remove it
		podman.DeleteSecret(SystemIDSecret, false)
		if err := GetSystemID(); err != nil {
			log.Warn().Err(err).Msg(L("Unable to fetch up to date systemid, using one from the provided configuration file"))
		}
	}

	authFile, cleaner, err := podman.PodmanLogin(hostData, flags.Registry, flags.SCC)
	if err != nil {
		return err
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

	return StartPod(systemd)
}

// Start the proxy services.
func StartPod(systemd podman.Systemd) error {
	ret := systemd.IsServiceRunning(podman.ProxyService)
	if ret {
		return systemd.RestartService(podman.ProxyService)
	}
	return systemd.EnableService(podman.ProxyService)
}

func getSystemIDEvent() ([]byte, error) {
	// Start the event listener in the background
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventListenerCmd := contextRunner(
		ctx,
		"venv-salt-call",
		"state.event",
		"tagmatch="+SystemIDEventResponse,
		"count=1",
		"--out=quiet",
	)
	var out bytes.Buffer

	log.Debug().Msg("Starting event listener")
	if err := eventListenerCmd.Std(&out).Start(); err != nil {
		return nil, err
	}

	// Allow event listener to start
	time.Sleep(time.Second)

	// Trigger the even
	fireEventCmd := newRunner(
		"venv-salt-call",
		"event.send",
		SystemIDEvent,
	)
	log.Debug().Msg("Asking for up to date systemid")
	if _, err := fireEventCmd.Exec(); err != nil {
		return nil, err
	}

	// Wait for the event listener to finish, we are waiting for one event at most 10s
	err := eventListenerCmd.Wait()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, ctx.Err()
		}
		return nil, err
	}
	return out.Bytes(), nil
}

func parseSystemIDEvent(event []byte) (string, error) {
	found := bytes.HasPrefix(event, []byte(SystemIDEventResponse))
	if !found {
		return "", errors.New(L("Not a system id event"))
	}
	jsonData := map[string]string{}
	err := json.Unmarshal(event[len(SystemIDEventResponse):], &jsonData)
	if err != nil {
		return "", err
	}
	data, ok := jsonData["data"]
	if !ok {
		return "", errors.New(L("System id not found in returned event"))
	}
	return data, nil
}

func GetSystemID() error {
	event, err := getSystemIDEvent()
	if err != nil {
		return err
	}

	systemid, err := parseSystemIDEvent(event)
	if err != nil {
		return err
	}
	log.Trace().Msgf("SystemID: %s", systemid)

	return podman.CreateSecret(SystemIDSecret, systemid)
}
