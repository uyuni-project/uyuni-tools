// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// GetExposedPorts returns the port exposed.
func GetExposedPorts(debug bool) []types.PortMap {
	ports := []types.PortMap{
		utils.NewPortMap("https", 443, 443),
		utils.NewPortMap("http", 80, 80),
	}
	ports = append(ports, utils.TCPPorts...)
	ports = append(ports, utils.TCPPodmanPorts...)
	ports = append(ports, utils.UDPPorts...)

	if debug {
		ports = append(ports, utils.DebugPorts...)
	}

	return ports
}

// GenerateServerSystemdService creates the server systemd service file.
func GenerateServerSystemdService(mirrorPath string, debug bool) error {
	ipv6Enabled := podman.HasIpv6Enabled(podman.UyuniNetwork)

	args := podman.GetCommonParams()

	if mirrorPath != "" {
		args = append(args, "-v", mirrorPath+":/mirror")
	}

	ports := GetExposedPorts(debug)
	if _, err := exec.LookPath("csp-billing-adapter"); err == nil {
		ports = append(ports, utils.NewPortMap("csp-billing", 18888, 18888))
		args = append(args, "-e ISPAYG=1")
	}

	data := templates.PodmanServiceTemplateData{
		Volumes:     utils.ServerVolumeMounts,
		NamePrefix:  "uyuni",
		Args:        strings.Join(args, " "),
		Ports:       ports,
		Network:     podman.UyuniNetwork,
		IPV6Enabled: ipv6Enabled,
	}
	if err := utils.WriteTemplateToFile(data, podman.GetServicePath("uyuni-server"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	return nil
}

// GenerateSystemdService creates a server systemd file.
func GenerateSystemdService(tz string, image string, debug bool, mirrorPath string, podmanArgs []string) error {
	err := podman.SetupNetwork(false)
	if err != nil {
		return utils.Errorf(err, L("cannot setup network"))
	}

	log.Info().Msg(L("Enabling system service"))
	if err := GenerateServerSystemdService(mirrorPath, debug); err != nil {
		return err
	}

	if err := podman.GenerateSystemdConfFile("uyuni-server", "generated.conf",
		"Environment=UYUNI_IMAGE="+image, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	config := fmt.Sprintf(`Environment=TZ=%s
Environment="PODMAN_EXTRA_ARGS=%s"
`, strings.TrimSpace(tz), strings.Join(podmanArgs, " "))

	if err := podman.GenerateSystemdConfFile("uyuni-server", "custom.conf", config, false); err != nil {
		return utils.Errorf(err, L("cannot generate systemd user configuration file"))
	}
	return podman.ReloadDaemon(false)
}

// UpdateSSLCertificate update SSL certificate.
func UpdateSSLCertificate(cnx *shared.Connection, chain *ssl.CaChain, serverPair *ssl.SSLPair) error {
	ssl.CheckPaths(chain, serverPair)

	// Copy the CAs, certificate and key to the container
	const certDir = "/tmp/uyuni-tools"
	if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "mkdir", "-p", certDir); err != nil {
		return errors.New(L("failed to create temporary folder on container to copy certificates to"))
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
		return utils.Errorf(err, L("cannot copy %s"), rootCaPath)
	}
	if err := cnx.Copy(serverPair.Cert, "server:"+serverCrtPath, "root", "root"); err != nil {
		return utils.Errorf(err, L("cannot copy %s"), serverCrtPath)
	}
	if err := cnx.Copy(serverPair.Key, "server:"+serverKeyPath, "root", "root"); err != nil {
		return utils.Errorf(err, L("cannot copy %s"), serverKeyPath)
	}

	for i, ca := range chain.Intermediate {
		caFilename := fmt.Sprintf("ca-%d.crt", i)
		caPath := path.Join(certDir, caFilename)
		args = append(args, "--intermediate-ca-file", caPath)
		if err := cnx.Copy(ca, "server:"+caPath, "root", "root"); err != nil {
			return utils.Errorf(err, L("cannot copy %s"), caPath)
		}
	}

	// Check and install then using mgr-ssl-cert-setup
	if out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", args...); err != nil {
		return utils.Errorf(err, L("failed to update SSL certificate: %s"), out)
	}

	// Clean the copied files and the now useless ssl-build
	if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "rm", "-rf", certDir); err != nil {
		return utils.Errorf(err, L("failed to remove copied certificate files in the container"))
	}

	const sslbuildPath = "/root/ssl-build"
	if cnx.TestExistenceInPod(sslbuildPath) {
		if err := utils.RunCmd("podman", "exec", podman.ServerContainerName, "rm", "-rf", sslbuildPath); err != nil {
			return utils.Errorf(err, L("failed to remove now useless ssl-build folder in the container"))
		}
	}

	// The services need to be restarted
	log.Info().Msg(L("Restarting services after updating the certificate"))
	if err := utils.RunCmd(
		"podman", "exec", podman.ServerContainerName, "systemctl", "restart", "postgresql.service",
	); err != nil {
		return err
	}
	return utils.RunCmdStdMapping(
		zerolog.DebugLevel, "podman", "exec", podman.ServerContainerName, "spacewalk-service", "restart",
	)
}

// RunMigration migrate an existing remote server to a container.
func RunMigration(
	preparedImage string,
	sshAuthSocket string,
	sshConfigPath string,
	sshKnownhostsPath string,
	sourceFqdn string,
	user string,
	prepare bool,
) (*utils.InspectResult, error) {
	scriptDir, cleaner, err := adm_utils.GenerateMigrationScript(sourceFqdn, user, false, prepare)
	if err != nil {
		return nil, utils.Errorf(err, L("cannot generate migration script"))
	}
	defer cleaner()

	extraArgs := []string{
		"--security-opt", "label=disable",
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

	log.Info().Msg(L("Migrating server"))
	if err := podman.RunContainer("uyuni-migration", preparedImage, utils.ServerVolumeMounts, extraArgs,
		[]string{"/var/lib/uyuni-tools/migrate.sh"}); err != nil {
		return nil, utils.Errorf(err, L("cannot run uyuni migration container"))
	}

	//now that everything is migrated, we need to fix SELinux permission
	for _, volumeMount := range utils.ServerVolumeMounts {
		mountPoint, err := GetMountPoint(volumeMount.Name)
		if err != nil {
			return nil, utils.Errorf(err, L("cannot inspect volume %s"), volumeMount)
		}
		if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "restorecon", "-F", "-r", "-v", mountPoint); err != nil {
			return nil, utils.Errorf(err, L("cannot restore %s SELinux permissions"), mountPoint)
		}
	}

	extractedData, err := utils.ReadInspectData[utils.InspectResult](path.Join(scriptDir, "data"))

	if err != nil {
		return nil, utils.Errorf(err, L("cannot read extracted data"))
	}

	return extractedData, nil
}

// RunPgsqlVersionUpgrade perform a PostgreSQL major upgrade.
func RunPgsqlVersionUpgrade(
	authFile string,
	registry string,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	oldPgsql string,
	newPgsql string,
) error {
	log.Info().Msgf(
		L("Previous PostgreSQL is %[1]s, new one is %[2]s. Performing a DB version upgrade…"), oldPgsql, newPgsql,
	)

	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()
	if newPgsql > oldPgsql {
		pgsqlVersionUpgradeContainer := "uyuni-upgrade-pgsql"
		extraArgs := []string{
			"-v", scriptDir + ":/var/lib/uyuni-tools/",
			"--security-opt", "label=disable",
		}

		upgradeImageURL := ""
		if upgradeImage.Name == "" {
			upgradeImageURL, err = utils.ComputeImage(registry, utils.DefaultTag, image,
				fmt.Sprintf("-migration-%s-%s", oldPgsql, newPgsql))
			if err != nil {
				return utils.Errorf(err, L("failed to compute image URL"))
			}
		} else {
			upgradeImageURL, err = utils.ComputeImage(registry, image.Tag, upgradeImage)
			if err != nil {
				return utils.Errorf(err, L("failed to compute image URL"))
			}
		}

		preparedImage, err := podman.PrepareImage(authFile, upgradeImageURL, image.PullPolicy, true)
		if err != nil {
			return err
		}

		log.Info().Msgf(L("Using database upgrade image %s"), preparedImage)

		pgsqlVersionUpgradeScriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(
			scriptDir, oldPgsql, newPgsql, false,
		)
		if err != nil {
			return utils.Errorf(err, L("cannot generate PostgreSQL database version upgrade script"))
		}

		err = podman.RunContainer(pgsqlVersionUpgradeContainer, preparedImage, utils.ServerVolumeMounts, extraArgs,
			[]string{"/var/lib/uyuni-tools/" + pgsqlVersionUpgradeScriptName})
		if err != nil {
			return err
		}
	}
	return nil
}

// RunPgsqlFinalizeScript run the script with all the action required to a db after upgrade.
func RunPgsqlFinalizeScript(serverImage string, schemaUpdateRequired bool, migration bool) error {
	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	extraArgs := []string{
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
		"--security-opt", "label=disable",
	}
	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"
	pgsqlFinalizeScriptName, err := adm_utils.GenerateFinalizePostgresScript(
		scriptDir, true, schemaUpdateRequired, true, migration, false,
	)
	if err != nil {
		return utils.Errorf(err, L("cannot generate PostgreSQL finalization script"))
	}
	err = podman.RunContainer(pgsqlFinalizeContainer, serverImage, utils.ServerVolumeMounts, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + pgsqlFinalizeScriptName})
	if err != nil {
		return err
	}
	return nil
}

// RunPostUpgradeScript run the script with the changes to apply after the upgrade.
func RunPostUpgradeScript(serverImage string) error {
	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()
	postUpgradeContainer := "uyuni-post-upgrade"
	extraArgs := []string{
		"-v", scriptDir + ":/var/lib/uyuni-tools/",
		"--security-opt", "label=disable",
	}
	postUpgradeScriptName, err := adm_utils.GeneratePostUpgradeScript(scriptDir, "localhost")
	if err != nil {
		return utils.Errorf(err, L("cannot generate PostgreSQL finalization script"))
	}
	err = podman.RunContainer(postUpgradeContainer, serverImage, utils.ServerVolumeMounts, extraArgs,
		[]string{"/var/lib/uyuni-tools/" + postUpgradeScriptName})
	if err != nil {
		return err
	}
	return nil
}

// Upgrade will upgrade server to the image given as attribute.
func Upgrade(
	authFile string,
	registry string,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	cocoFlags adm_utils.CocoFlags,
	hubXmlrpcFlags adm_utils.HubXmlrpcFlags,
) error {
	if err := CallCloudGuestRegistryAuth(); err != nil {
		return err
	}

	serverImage, err := utils.ComputeImage(registry, utils.DefaultTag, image)
	if err != nil {
		return errors.New(L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, serverImage, image.PullPolicy, true)
	if err != nil {
		return err
	}

	inspectedValues, err := Inspect(preparedImage)
	if err != nil {
		return utils.Errorf(err, L("cannot inspect podman values"))
	}

	cnx := shared.NewConnection("podman", podman.ServerContainerName, "")

	if err := adm_utils.SanityCheck(cnx, inspectedValues, preparedImage); err != nil {
		return err
	}

	if err := podman.StopService(podman.ServerService); err != nil {
		return utils.Errorf(err, L("cannot stop service"))
	}

	defer func() {
		err = podman.StartService(podman.ServerService)
	}()
	if inspectedValues.ImagePgVersion > inspectedValues.CurrentPgVersion {
		log.Info().Msgf(
			L("Previous postgresql is %[1]s, instead new one is %[2]s. Performing a DB version upgrade…"),
			inspectedValues.CurrentPgVersion, inspectedValues.ImagePgVersion,
		)
		if err := RunPgsqlVersionUpgrade(
			authFile, registry, image, upgradeImage, inspectedValues.CurrentPgVersion, inspectedValues.ImagePgVersion,
		); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	} else if inspectedValues.ImagePgVersion == inspectedValues.CurrentPgVersion {
		log.Info().Msgf(L("Upgrading to %s without changing PostgreSQL version"), inspectedValues.UyuniRelease)
	} else {
		return fmt.Errorf(
			L("trying to downgrade PostgreSQL from %[1]s to %[2]s"),
			inspectedValues.CurrentPgVersion, inspectedValues.ImagePgVersion,
		)
	}

	schemaUpdateRequired := inspectedValues.CurrentPgVersion != inspectedValues.ImagePgVersion
	if err := RunPgsqlFinalizeScript(preparedImage, schemaUpdateRequired, false); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL finalize script"))
	}

	if err := RunPostUpgradeScript(preparedImage); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	if err := podman.CleanSystemdConfFile("uyuni-server"); err != nil {
		return err
	}

	if err := podman.GenerateSystemdConfFile("uyuni-server", "generated.conf",
		"Environment=UYUNI_IMAGE="+preparedImage, true,
	); err != nil {
		return err
	}

	if err := podman.ReloadDaemon(false); err != nil {
		return err
	}

	if err := updateServerSystemdService(); err != nil {
		return err
	}
	log.Info().Msg(L("Waiting for the server to start…"))

	err = coco.Upgrade(authFile, registry, cocoFlags, image,
		inspectedValues.DBPort, inspectedValues.DBName, inspectedValues.DBUser, inspectedValues.DBPassword)
	if err != nil {
		return utils.Errorf(err, L("error upgrading confidential computing service."))
	}

	if err := hub.Upgrade(
		authFile, registry, image.PullPolicy, image.Tag, hubXmlrpcFlags,
	); err != nil {
		return err
	}

	return podman.ReloadDaemon(false)
}

var runCmdOutput = utils.RunCmdOutput

func hasDebugPorts(definition []byte) bool {
	return regexp.MustCompile(`-p 8003:8003`).Match(definition)
}

func getMirrorPath(definition []byte) string {
	mirrorPath := ""
	finder := regexp.MustCompile(`-v +([^:]+):/mirror[[:space:]]`)
	submatches := finder.FindStringSubmatch(string(definition))
	if len(submatches) == 2 {
		mirrorPath = submatches[1]
	}
	return mirrorPath
}

func updateServerSystemdService() error {
	out, err := runCmdOutput(zerolog.DebugLevel, "systemctl", "cat", podman.ServerService)
	if err != nil {
		return utils.Errorf(err, "failed to get %s systemd service definition", podman.ServerService)
	}

	return GenerateServerSystemdService(getMirrorPath(out), hasDebugPorts(out))
}

// Inspect check values on a given image and deploy.
func Inspect(preparedImage string) (*utils.ServerInspectData, error) {
	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return nil, err
	}
	defer cleaner()

	inspector := utils.NewServerInspector(scriptDir)
	if err := inspector.GenerateScript(); err != nil {
		return nil, err
	}

	podmanArgs := []string{
		"-v", scriptDir + ":" + utils.InspectContainerDirectory,
		"--security-opt", "label=disable",
	}

	err = podman.RunContainer("uyuni-inspect", preparedImage, utils.ServerVolumeMounts, podmanArgs,
		[]string{utils.InspectContainerDirectory + "/" + utils.InspectScriptFilename})
	if err != nil {
		return nil, err
	}

	inspectResult, err := inspector.ReadInspectData()
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect data"))
	}

	return inspectResult, err
}

// CallCloudGuestRegistryAuth calls cloudguestregistryauth if it is available.
func CallCloudGuestRegistryAuth() error {
	cloudguestregistryauth := "cloudguestregistryauth"

	path, err := exec.LookPath(cloudguestregistryauth)
	if err == nil {
		// the binary is installed
		return utils.RunCmdStdMapping(zerolog.DebugLevel, path)
	}
	// silently ignore error if it is missing
	return nil
}

// GetMountPoint return folder where a given volume is mounted.
func GetMountPoint(volumeName string) (string, error) {
	args := []string{"volume", "inspect", "--format", "{{.Mountpoint}}", volumeName}
	mountPoint, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(mountPoint), "\n"), nil
}
