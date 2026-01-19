// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/pgsql"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/saline"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.NewSystemd()

// GetExposedPorts returns the port exposed.
func GetExposedPorts(debug bool) []types.PortMap {
	ports := utils.GetServerPorts(debug)
	ports = append(ports, utils.NewPortMap(utils.WebServiceName, "https", 443, 443))
	ports = append(ports, utils.TCPPodmanPorts...)
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
		ports = append(ports, utils.NewPortMap("csp", "csp-billing", 18888, 18888))
		args = append(args, "-e ISPAYG=1")
	}

	data := templates.PodmanServiceTemplateData{
		Volumes:     utils.ServerVolumeMounts,
		NamePrefix:  "uyuni",
		Args:        strings.Join(args, " "),
		Ports:       ports,
		Network:     podman.UyuniNetwork,
		IPV6Enabled: ipv6Enabled,
		CaSecret:    podman.CASecret,
		CaPath:      ssl.CAContainerPath,
		CertSecret:  podman.SSLCertSecret,
		CertPath:    ssl.ServerCertPath,
		KeySecret:   podman.SSLKeySecret,
		KeyPath:     ssl.ServerCertKeyPath,
		DBCaSecret:  podman.DBCASecret,
		DBCaPath:    ssl.DBCAContainerPath,
	}
	if err := utils.WriteTemplateToFile(data, podman.GetServicePath("uyuni-server"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	return nil
}

// GenerateSystemdService creates a server systemd file.
func GenerateSystemdService(
	systemd podman.Systemd,
	tz string,
	image string,
	debug bool,
	mirrorPath string,
	podmanArgs []string,
) error {
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

	if err := podman.GenerateSystemdConfFile("uyuni-server", podman.CustomConf, config, false); err != nil {
		return utils.Errorf(err, L("cannot generate systemd user configuration file"))
	}
	return systemd.ReloadDaemon(false)
}

func RunSSLMigration(
	preparedImage string,
	sshAuthSocket string,
	sshConfigPath string,
	sshKnownhostsPath string,
	sourceFqdn string,
	user string,
) (*utils.InspectResult, error) {
	scriptDir, cleaner, err := utils.TempDir()
	defer cleaner()
	if err != nil {
		return nil, err
	}

	t := templates.SSLMigrateScriptTemplateData{
		Volumes:    utils.SSLMigrationVolumeMounts,
		SourceFqdn: sourceFqdn,
		User:       user,
	}

	scriptBuilder := new(strings.Builder)
	if err := t.Render(scriptBuilder); err != nil {
		return nil, utils.Error(err, L("failed to generate SSL migration script"))
	}

	script := scriptBuilder.String()

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

	log.Info().Msgf(L("Migrating SSL certificates from the source server %s"), sourceFqdn)
	if err := podman.RunContainer("uyuni-ssl-migration", preparedImage, utils.SSLMigrationVolumeMounts, extraArgs,
		[]string{"bash", "-e", "-c", script}); err != nil {
		return nil, utils.Errorf(err, L("cannot run uyuni SSL migration container"))
	}

	// now that everything is migrated, we need to fix SELinux permission
	if err := restoreSELinuxContext(utils.SSLMigrationVolumeMounts); err != nil {
		return nil, err
	}

	dataPath := path.Join(scriptDir, "data")
	data, err := os.ReadFile(dataPath)
	if err != nil {
		log.Fatal().Err(err).Msgf(L("Failed to read file %s"), dataPath)
	}

	extractedData, err := utils.ReadInspectData[utils.InspectResult](data)

	if err != nil {
		return nil, utils.Errorf(err, L("cannot read extracted data"))
	}

	return extractedData, nil
}

func restoreSELinuxContext(volumes []types.VolumeMount) error {
	if utils.IsInstalled("restorecon") {
		for _, volumeMount := range volumes {
			mountPoint, err := GetMountPoint(volumeMount.Name)
			if err != nil {
				return utils.Errorf(err, L("cannot inspect volume %s"), volumeMount)
			}
			if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "restorecon", "-F", "-r", "-v", mountPoint); err != nil {
				return utils.Errorf(err, L("cannot restore %s SELinux permissions"), mountPoint)
			}
		}
	}
	return nil
}

var prepareImage = podman.PrepareImage
var runContainer = podman.RunContainer

// RunPgsqlVersionUpgrade perform a PostgreSQL major upgrade.
func RunPgsqlVersionUpgrade(
	authFile string,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	volumeMounts []types.VolumeMount,
) error {
	pgsqlVersionUpgradeContainer := "uyuni-upgrade-pgsql"
	extraArgs := []string{
		"--security-opt", "label=disable",
		"--tmpfs", "/tmp:rw,mode=1777",
	}

	if podman.HasSecret(podman.DBCASecret) {
		extraArgs = append(extraArgs,
			"--secret", fmt.Sprintf("%s,type=mount,target=%s", podman.DBCASecret, ssl.DBCAContainerPath),
		)
	}

	if podman.HasSecret(podman.DBSSLKeySecret) {
		extraArgs = append(extraArgs,
			"--secret", fmt.Sprintf("%s,type=mount,uid=999,mode=0400,target=%s", podman.DBSSLKeySecret, ssl.DBCertKeyPath),
		)
	}
	if podman.HasSecret(podman.DBSSLCertSecret) {
		extraArgs = append(extraArgs,
			"--secret", fmt.Sprintf("%s,type=mount,target=%s", podman.DBSSLCertSecret, ssl.DBCertPath),
		)
	}

	upgradeImageURL, err := utils.ComputeImage(image.Registry.Host, image.Tag, upgradeImage)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := prepareImage(authFile, upgradeImageURL, image.PullPolicy, true)
	if err != nil {
		return err
	}

	log.Info().Msgf(L("Using database upgrade image %s"), preparedImage)

	return runContainer(pgsqlVersionUpgradeContainer, preparedImage, volumeMounts, extraArgs,
		[]string{})
}

// RunPgsqlFinalizeScript run the script with all the action required to a db after upgrade.
func RunPgsqlFinalizeScript(serverImage string, schemaUpdateRequired bool, collationChange bool) error {
	if !schemaUpdateRequired && !collationChange {
		log.Info().Msg(L("No need to run database finalization script"))
		return nil
	}

	env := map[string]string{
		"RUN_REINDEX":       strconv.FormatBool(collationChange),
		"RUN_SCHEMA_UPDATE": strconv.FormatBool(schemaUpdateRequired),
	}

	extraArgs := []string{
		"--security-opt", "label=disable",
		"--network", podman.UyuniNetwork,
	}

	for key, value := range env {
		extraArgs = append(extraArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"

	return podman.RunContainer(pgsqlFinalizeContainer, serverImage, utils.ServerVolumeMounts, extraArgs,
		[]string{"/usr/bin/sh", "-e", "-c", "/docker-entrypoint-init.d/90-pgsqlFinalize.sh"})
}

// RunPostUpgradeScript run the script with the changes to apply after the upgrade.
func RunPostUpgradeScript(serverImage string) error {
	postUpgradeContainer := "uyuni-post-upgrade"
	extraArgs := []string{
		"--security-opt", "label=disable",
	}
	script, err := adm_utils.GeneratePostUpgradeScript()
	if err != nil {
		return utils.Errorf(err, L("cannot generate PostgreSQL finalization script"))
	}
	// Post upgrade script expects some commands to fail and checks their result, don't use sh -e.
	return podman.RunContainer(postUpgradeContainer, serverImage, utils.ServerVolumeMounts, extraArgs,
		[]string{"bash", "-c", script})
}

// Upgrade will upgrade server to the image given as attribute.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	db adm_utils.DBFlags,
	reportdb adm_utils.DBFlags,
	ssl adm_utils.InstallSSLFlags,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	cocoFlags adm_utils.CocoFlags,
	hubXmlrpcFlags adm_utils.HubXmlrpcFlags,
	salineFlags adm_utils.SalineFlags,
	pgsqlFlags types.PgsqlFlags,
	tz string,
) error {
	// Calling cloudguestregistryauth only makes sense if using the cloud provider registry.
	// This check assumes users won't use custom registries that are not the cloud provider one on a cloud image.
	if !strings.HasPrefix(image.Registry.Host, "registry.suse.com") {
		if err := CallCloudGuestRegistryAuth(); err != nil {
			return err
		}
	}

	// Prepare Uyuni network, migration container needs to run in the same network as resulting image
	err := podman.SetupNetwork(false)
	if err != nil {
		return utils.Errorf(err, L("cannot setup network"))
	}

	fqdn, err := utils.GetFqdn([]string{})
	if err != nil {
		return err
	}

	preparedServerImage, preparedPgsqlImage, err := podman.PrepareImages(authFile, image, pgsqlFlags)
	if err != nil {
		return utils.Errorf(err, L("cannot prepare images"))
	}

	inspectedValues, err := prepareHost(preparedServerImage, preparedPgsqlImage)
	if err != nil {
		return err
	}

	if systemd.HasService(podman.ServerService) {
		if err := systemd.StopService(podman.ServerService); err != nil {
			return utils.Errorf(err, L("cannot stop service"))
		}
		defer func() {
			err = systemd.StartService(podman.ServerService)
		}()
	}
	if systemd.HasService(podman.DBService) {
		if err := systemd.StopService(podman.DBService); err != nil {
			return utils.Errorf(err, L("cannot stop service"))
		}
		defer func() {
			err = systemd.StartService(podman.DBService)
		}()
	}

	oldPgVersion, _ := strconv.Atoi(inspectedValues.ContainerInspectData.PgVersion)
	newPgVersion, _ := strconv.Atoi(inspectedValues.DBInspectData.PgVersion)

	if newPgVersion > oldPgVersion {
		log.Info().Msgf(L("Initiating PostgreSQL upgrade from version %[1]d to %[2]d"), oldPgVersion, newPgVersion)

		pgsqlMountpoint, err := podman.GetVolumeMountPoint(utils.VarPgsqlDataVolumeMount.Name)
		if err != nil {
			return utils.Errorf(err, L("cannot find volume %s"), utils.VarPgsqlDataVolumeMount.Name)
		}

		targetPath := path.Join(pgsqlMountpoint, "..", "_data")
		upgradeVolumeMounts := []types.VolumeMount{
			{
				MountPath: "/migration/target",
				Name:      targetPath,
			},
			utils.EtcTLSTmpVolumeMount,
		}

		backupPath := path.Join(pgsqlMountpoint, "..", "_data_old")

		if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "mv", targetPath, backupPath); err != nil {
			return utils.Errorf(err, L("cannot move %s"), targetPath)
		}

		if strings.HasPrefix(inspectedValues.ContainerInspectData.SuseManagerRelease, "5.0") {
			upgradeVolumeMounts = append(upgradeVolumeMounts, types.VolumeMount{
				MountPath: "/migration/source",
				Name:      path.Join(backupPath, "data"),
			})
		} else {
			upgradeVolumeMounts = append(upgradeVolumeMounts, types.VolumeMount{
				MountPath: "/migration/source",
				Name:      backupPath,
			})
		}

		if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "mkdir", "-p", targetPath); err != nil {
			return utils.Errorf(err, L("cannot mkdir %s"), targetPath)
		}

		log.Warn().Msg(L("Data will be copied during this process. Please ensure sufficient disk space is available."))
		if err := RunPgsqlVersionUpgrade(authFile, image, upgradeImage, upgradeVolumeMounts); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	} else if newPgVersion == oldPgVersion {
		log.Info().Msg(L("Upgrading without changing PostgreSQL version"))
	} else {
		return fmt.Errorf(
			L("trying to downgrade PostgreSQL from %[1]d to %[2]d"),
			oldPgVersion, newPgVersion,
		)
	}

	if inspectedValues.DBHost == "localhost" ||
		inspectedValues.ReportDBHost == "localhost" {
		log.Info().Msgf(L("Configuring split PostgreSQL container"))

		if err := PrepareSSLCertificates(preparedServerImage, &ssl, tz, fqdn); err != nil {
			return err
		}
	}

	if err := configureDBContainer(
		preparedServerImage, preparedPgsqlImage, systemd, db, reportdb); err != nil {
		return utils.Errorf(err, L("cannot configure db container"))
	}

	if err := pgsql.Upgrade(preparedPgsqlImage, systemd); err != nil {
		return err
	}

	schemaUpdateRequired := oldPgVersion != newPgVersion
	collationChange := inspectedValues.ServerInspectData.LibcVersion != inspectedValues.ContainerInspectData.LibcVersion
	if err := RunPgsqlFinalizeScript(preparedServerImage, schemaUpdateRequired, collationChange); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL finalize script"))
	}

	if err := RunPostUpgradeScript(preparedServerImage); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	if err := podman.CleanSystemdConfFile("uyuni-server"); err != nil {
		return err
	}

	if err := podman.GenerateSystemdConfFile("uyuni-server", "generated.conf",
		"Environment=UYUNI_IMAGE="+preparedServerImage, true,
	); err != nil {
		return err
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	if err := UpdateServerSystemdService(); err != nil {
		return err
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	log.Info().Msg(L("Waiting for the server to start…"))
	cnx := shared.NewConnection("podman", podman.ServerContainerName, "")
	if err := systemd.StartService(podman.ServerService); err != nil {
		return utils.Error(err, L("cannot start service"))
	}

	if err := cnx.WaitForHealthcheck(); err != nil {
		log.Warn().Err(err)
	}

	inspectedDB := adm_utils.DBFlags{
		Name:     inspectedValues.DBName,
		Port:     inspectedValues.DBPort,
		User:     inspectedValues.DBUser,
		Password: inspectedValues.DBPassword,
		Host:     db.Host,
	}

	err = coco.Upgrade(systemd, authFile, cocoFlags, image, inspectedDB)

	if err != nil {
		return utils.Errorf(err, L("error upgrading confidential computing service."))
	}

	if err := hub.Upgrade(
		systemd, authFile, image, hubXmlrpcFlags,
	); err != nil {
		return err
	}

	if err := saline.Upgrade(systemd, authFile, image, salineFlags, utils.GetLocalTimezone()); err != nil {
		return utils.Errorf(err, L("error upgrading saline service."))
	}

	return systemd.ReloadDaemon(false)
}

func WaitForSystemStart(
	systemd podman.Systemd,
	cnx *shared.Connection,
	image string,
	tz string,
	debug bool,
	mirrorPath string,
	podmanArgs []string,
) error {
	err := GenerateSystemdService(
		systemd, tz, image, debug, mirrorPath, podmanArgs,
	)
	if err != nil {
		return err
	}

	log.Info().Msg(L("Waiting for the server to start…"))
	if err := systemd.EnableService(podman.ServerService); err != nil {
		return utils.Error(err, L("cannot enable service"))
	}

	return cnx.WaitForHealthcheck()
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

// UpdateServerSystemdService refreshes the server systemd service file.
func UpdateServerSystemdService() error {
	out, err := runCmdOutput(zerolog.DebugLevel, "systemctl", "cat", podman.ServerService)
	if err != nil {
		return utils.Errorf(err, "failed to get %s systemd service definition", podman.ServerService)
	}

	return GenerateServerSystemdService(getMirrorPath(out), hasDebugPorts(out))
}

// RunSplitContainerSettings migrate to separate postgres container.
func RunSplitContainerSettings(serverImage string, dbHost string, reportDBHost string) error {
	data := templates.SplitContainerSettingsScriptTemplateData{
		DBHost:       dbHost,
		ReportDBHost: reportDBHost,
	}

	scriptBuilder := new(strings.Builder)
	if err := data.Render(scriptBuilder); err != nil {
		return utils.Error(err, L("failed to generate postgresql migration script"))
	}

	podmanArgs := []string{
		"--security-opt", "label=disable",
	}
	return podman.RunContainer("uyuni-db-migrate", serverImage, utils.DatabaseMigrationVolumeMounts, podmanArgs,
		[]string{"bash", "-e", "-c", scriptBuilder.String()})
}

// RunConfigPgsl setup postgres container.
func RunConfigPgsl(pgsqlImage string) error {
	podmanArgs := []string{
		"--security-opt", "label=disable",
		"--secret", fmt.Sprintf("%s,type=mount,target=%s", podman.DBCASecret, ssl.DBCAContainerPath),
		"--secret", fmt.Sprintf("%s,type=mount,uid=999,mode=0400,target=%s", podman.DBSSLKeySecret, ssl.DBCertKeyPath),
		"--secret", fmt.Sprintf("%s,type=mount,target=%s", podman.DBSSLCertSecret, ssl.DBCertPath),
		"--entrypoint", "/docker-entrypoint-initdb.d/uyuni-postgres-config.sh",
	}
	if err := podman.RunContainer("uyuni-db-config", pgsqlImage, utils.PgsqlRequiredVolumeMounts,
		podmanArgs, []string{}); err != nil {
		return err
	}
	return systemd.RestartService(podman.DBService)
}

// CallCloudGuestRegistryAuth calls cloudguestregistryauth if it is available.
func CallCloudGuestRegistryAuth() error {
	cloudguestregistryauth := "cloudguestregistryauth"

	path, err := exec.LookPath(cloudguestregistryauth)
	if err == nil {
		if err := utils.RunCmdStdMapping(zerolog.DebugLevel, path); err != nil && isPAYG() {
			// Not being registered against the cloud registry is  not an error on BYOS.
			return err
		} else if err != nil {
			log.Info().Msg(L("The above error is only relevant if using a public cloud provider registry"))
		}
	}
	// silently ignore error if it is missing
	return nil
}

func isPAYG() bool {
	flavorCheckPath := "/usr/bin/instance-flavor-check"
	if utils.FileExists(flavorCheckPath) {
		out, _ := utils.RunCmdOutput(zerolog.DebugLevel, flavorCheckPath)
		return strings.TrimSpace(string(out)) == "PAYG"
	}
	return false
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

// GetSSHAuthSocket returns the SSH_AUTH_SOCK environment variable value.
func GetSSHAuthSocket() string {
	path := os.Getenv("SSH_AUTH_SOCK")
	if len(path) == 0 {
		log.Fatal().Msg(L("SSH_AUTH_SOCK is not defined, start an SSH agent and try again"))
	}
	return path
}

// GetSSHPaths returns the user SSH config and known_hosts paths.
func GetSSHPaths() (string, string) {
	// Find ssh config to mount it in the container
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg(L("Failed to find home directory to look for SSH config"))
	}
	sshConfigPath := filepath.Join(homedir, ".ssh", "config")
	sshKnownhostsPath := filepath.Join(homedir, ".ssh", "known_hosts")

	if !utils.FileExists(sshConfigPath) {
		sshConfigPath = ""
	}

	if !utils.FileExists(sshKnownhostsPath) {
		sshKnownhostsPath = ""
	}

	return sshConfigPath, sshKnownhostsPath
}

func prepareHost(
	preparedServerImage string,
	preparedPgsqlImage string,
) (*utils.InspectData, error) {
	inspectedValues, err := podman.Inspect(preparedServerImage, preparedPgsqlImage)
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect podman values"))
	}

	return inspectedValues, adm_utils.SanityCheck(inspectedValues)
}

func configureDBContainer(
	serverImage string,
	pgsqlImage string,
	systemd podman.Systemd,
	db adm_utils.DBFlags,
	reportdb adm_utils.DBFlags,
) error {
	if err := RunSplitContainerSettings(serverImage, "db", "reportdb"); err != nil {
		return utils.Errorf(err, L("PostgreSQL migration failure"))
	}

	// Create all the database credentials secrets
	if err := podman.CreateCredentialsSecretsIfMissing(
		podman.DBUserSecret, db.User,
		podman.DBPassSecret, db.Password,
	); err != nil {
		return err
	}

	if err := podman.CreateCredentialsSecretsIfMissing(
		podman.ReportDBUserSecret, reportdb.User,
		podman.ReportDBPassSecret, reportdb.Password,
	); err != nil {
		return err
	}

	if db.IsLocal() {
		if !podman.HasSecret(podman.DBAdminUserSecret) && !podman.HasSecret(podman.DBAdminPassSecret) {
			// The admin password is not needed for external databases
			if err := podman.CreateCredentialsSecrets(
				podman.DBAdminUserSecret, db.Admin.User,
				podman.DBAdminPassSecret, db.Admin.Password,
			); err != nil {
				return err
			}
		} else {
			log.Info().Msgf(
				L("Secrets %[1]s and %[2]s, already exists"), podman.DBAdminUserSecret, podman.DBAdminPassSecret)
		}

		// Run the DB container setup if the user doesn't set a custom host name for it.
		if err := pgsql.SetupPgsql(systemd, pgsqlImage); err != nil {
			return err
		}
	} else {
		log.Info().Msgf(
			L("Skipped database container setup to use external database %s"),
			db.Host,
		)
	}
	return RunConfigPgsl(pgsqlImage)
}
