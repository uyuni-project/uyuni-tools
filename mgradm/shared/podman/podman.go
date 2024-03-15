// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

// GetCommonParams splits the common arguments.
func GetCommonParams() []string {
	return strings.Split(commonArgs, " ")
}

// GetExposedPorts returns the port exposed.
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

// GenerateSystemdService creates a serverY systemd file.
func GenerateSystemdService(tz string, image string, debug bool, podmanArgs []string) error {
	if err := podman.SetupNetwork(); err != nil {
		return fmt.Errorf("cannot setup network: %s", err)
	}

	log.Info().Msg("Enabling system service")
	data := templates.PodmanServiceTemplateData{
		Volumes:    utils.ServerVolumeMounts,
		NamePrefix: "uyuni",
		Args:       commonArgs + " " + strings.Join(podmanArgs, " "),
		Ports:      GetExposedPorts(debug),
		Timezone:   tz,
		Network:    podman.UyuniNetwork,
	}
	if err := utils.WriteTemplateToFile(data, podman.GetServicePath("uyuni-server"), 0555, false); err != nil {
		return fmt.Errorf("failed to generate systemd service unit file: %s", err)
	}

	if err := podman.GenerateSystemdConfFile("uyuni-server", "Service", "Environment=UYUNI_IMAGE="+image); err != nil {
		return fmt.Errorf("cannot generate systemd conf file: %s", err)
	}
	return podman.ReloadDaemon(false)
}

// UpdateSslCertificate update SSL certificate.
func UpdateSslCertificate(cnx *shared.Connection, chain *ssl.CaChain, serverPair *ssl.SslPair) error {
	ssl.CheckPaths(chain, serverPair)

	// Copy the CAs, certificate and key to the container
	const certDir = "/tmp/uyuni-tools"
	if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "mkdir", "-p", certDir); err != nil {
		return fmt.Errorf("failed to create temporary folder on container to copy certificates to")
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

	if err := cnx.Copy(chain.Root, "server:"+rootCaPath, "root", "root"); err != nil {
		return fmt.Errorf("cannot copy %s: %s", rootCaPath, err)
	}
	if err := cnx.Copy(serverPair.Cert, "server:"+serverCrtPath, "root", "root"); err != nil {
		return fmt.Errorf("cannot copy %s: %s", serverCrtPath, err)
	}
	if err := cnx.Copy(serverPair.Key, "server:"+serverKeyPath, "root", "root"); err != nil {
		return fmt.Errorf("cannot copy %s: %s", serverKeyPath, err)
	}

	for i, ca := range chain.Intermediate {
		caFilename := fmt.Sprintf("ca-%d.crt", i)
		caPath := path.Join(certDir, caFilename)
		args = append(args, "--intermediate-ca-file", caPath)
		if err := cnx.Copy(ca, "server:"+caPath, "root", "root"); err != nil {
			return fmt.Errorf("cannot copy %s: %s", caPath, err)
		}
	}

	// Check and install then using mgr-ssl-cert-setup
	if _, err := utils.RunCmdOutput(zerolog.InfoLevel, "podman", args...); err != nil {
		return fmt.Errorf("failed to update SSL certificate")
	}

	// Clean the copied files and the now useless ssl-build
	if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "rm", "-rf", certDir); err != nil {
		return fmt.Errorf("failed to remove copied certificate files in the container")
	}

	const sslbuildPath = "/root/ssl-build"
	if cnx.TestExistenceInPod(sslbuildPath) {
		if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "rm", "-rf", sslbuildPath); err != nil {
			return fmt.Errorf("failed to remove now useless ssl-build folder in the container")
		}
	}

	// The services need to be restarted
	log.Info().Msg("Restarting services after updating the certificate")
	return utils.RunCmdStdMapping(zerolog.DebugLevel, "podman", "exec", podman.ServerContainerName, "spacewalk-service", "restart")
}

// RunContainer execute a container.
func RunContainer(name string, image string, extraArgs []string, cmd []string) error {
	podmanArgs := append([]string{"run", "--name", name}, GetCommonParams()...)
	podmanArgs = append(podmanArgs, extraArgs...)
	for _, volume := range utils.ServerVolumeMounts {
		podmanArgs = append(podmanArgs, "-v", volume.Name+":"+volume.MountPath)
	}
	podmanArgs = append(podmanArgs, image)
	podmanArgs = append(podmanArgs, cmd...)

	err := utils.RunCmdStdMapping(zerolog.DebugLevel, "podman", podmanArgs...)
	if err != nil {
		return fmt.Errorf("failed to run %s container: %s", name, err)
	}

	return nil
}

// RunMigration migrate an existing remote server to a container.
func RunMigration(serverImage string, pullPolicy string, sshAuthSocket string, sshConfigPath string, sshKnownhostsPath string, sourceFqdn string) (string, string, string, error) {
	scriptDir, err := adm_utils.GenerateMigrationScript(sourceFqdn, false)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot generate migration script: %s", err)
	}
	defer os.RemoveAll(scriptDir)

	extraArgs := []string{
		"--security-opt", "label:disable",
		"-e", "SSH_AUTH_SOCK",
		"-v", filepath.Dir(sshAuthSocket) + ":" + filepath.Dir(sshAuthSocket),
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
	}

	if sshConfigPath != "" {
		extraArgs = append(extraArgs, "-v", sshConfigPath+":/tmp/ssh_config")
	}

	if sshKnownhostsPath != "" {
		extraArgs = append(extraArgs, "-v", sshKnownhostsPath+":/etc/ssh/ssh_known_hosts")
	}

	inspectedHostValues, err := adm_utils.InspectHost()
	if err != nil {
		return "", "", "", fmt.Errorf("cannot inspect host values: %s", err)
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	err = podman.PrepareImage(serverImage, pullPolicy, pullArgs...)
	if err != nil {
		return "", "", "", err
	}

	log.Info().Msg("Migrating server")
	if err := RunContainer("uyuni-migration", serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/migrate.sh"}); err != nil {
		return "", "", "", fmt.Errorf("cannot run uyuni migration container: %s", err)
	}
	tz, oldPgVersion, newPgVersion, err := adm_utils.ReadContainerData(scriptDir)

	if err != nil {
		return "", "", "", fmt.Errorf("cannot read extracted data: %s", err)
	}

	return tz, oldPgVersion, newPgVersion, nil
}

// RunPgsqlVersionUpgrade perform a PostgreSQL major upgrade.
func RunPgsqlVersionUpgrade(image types.ImageFlags, migrationImage types.ImageFlags, oldPgsql string, newPgsql string) error {
	log.Info().Msgf("Previous postgresql is %s, instead new one is %s. Performing a DB version upgrade...", oldPgsql, newPgsql)

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	if newPgsql > oldPgsql {
		pgsqlVersionUpgradeContainer := "uyuni-upgrade-pgsql"
		extraArgs := []string{
			"-v", scriptDir + ":/var/lib/uyuni-tools/",
			"--security-opt", "label:disable",
		}

		migrationImageUrl := ""
		if migrationImage.Name == "" {
			migrationImageUrl, err = utils.ComputeImage(image.Name, image.Tag, fmt.Sprintf("-migration-%s-%s", oldPgsql, newPgsql))
			if err != nil {
				return fmt.Errorf("failed to compute image URL %s", err)
			}
		} else {
			migrationImageUrl, err = utils.ComputeImage(migrationImage.Name, image.Tag)
			if err != nil {
				return fmt.Errorf("failed to compute image URL %s", err)
			}
		}

		inspectedHostValues, err := adm_utils.InspectHost()
		if err != nil {
			return fmt.Errorf("cannot inspect host values: %s", err)
		}

		pullArgs := []string{}
		_, scc_user_exist := inspectedHostValues["host_scc_username"]
		_, scc_user_password := inspectedHostValues["host_scc_password"]
		if scc_user_exist && scc_user_password {
			pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
		}

		err = podman.PrepareImage(migrationImageUrl, image.PullPolicy, pullArgs...)
		if err != nil {
			return err
		}

		log.Info().Msgf("Using migration image %s", migrationImageUrl)

		pgsqlVersionUpgradeScriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(scriptDir, oldPgsql, newPgsql, false)
		if err != nil {
			return fmt.Errorf("cannot generate postgresql database version upgrade script %s", err)
		}

		err = RunContainer(pgsqlVersionUpgradeContainer, migrationImageUrl, extraArgs,
			[]string{"/var/lib/uyuni-tools/" + pgsqlVersionUpgradeScriptName})
		if err != nil {
			return err
		}
	}
	return nil
}

// RunPgsqlFinalizeScript run the script with all the action required to a db after upgrade.
func RunPgsqlFinalizeScript(serverImage string, schemaUpdateRequired bool) error {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}

	extraArgs := []string{
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
		"--security-opt", "label:disable",
	}
	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"
	pgsqlFinalizeScriptName, err := adm_utils.GenerateFinalizePostgresScript(scriptDir, true, schemaUpdateRequired, true, true, false)
	if err != nil {
		return fmt.Errorf("cannot generate postgresql finalization script %s", err)
	}
	err = RunContainer(pgsqlFinalizeContainer, serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + pgsqlFinalizeScriptName})
	if err != nil {
		return err
	}
	return nil
}

// RunPostUpgradeScript run the script with the changes to apply after the upgrade.
func RunPostUpgradeScript(serverImage string) error {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	postUpgradeContainer := "uyuni-post-upgrade"
	extraArgs := []string{
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
		"--security-opt", "label:disable",
	}
	postUpgradeScriptName, err := adm_utils.GeneratePostUpgradeScript(scriptDir, "localhost")
	if err != nil {
		return fmt.Errorf("cannot generate postgresql finalization script %s", err)
	}
	err = RunContainer(postUpgradeContainer, serverImage, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + postUpgradeScriptName})
	if err != nil {
		return err
	}
	return nil
}
