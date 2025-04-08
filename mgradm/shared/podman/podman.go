// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	"github.com/spf13/viper"
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

var systemd podman.Systemd = podman.SystemdImpl{}

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

	if err := podman.GenerateSystemdConfFile("uyuni-server", "custom.conf", config, false); err != nil {
		return utils.Errorf(err, L("cannot generate systemd user configuration file"))
	}
	return systemd.ReloadDaemon(false)
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
	scriptDir, cleaner, err := adm_utils.GenerateMigrationScript(
		sourceFqdn,
		user,
		false,
		prepare,
		"uyuni-pgsql-server.mgr.internal",
		"uyuni-pgsql-server.mgr.internal",
	)
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
	if err := podman.RunContainer("uyuni-migration", preparedImage, utils.ServerMigrationVolumeMounts, extraArgs,
		[]string{"/var/lib/uyuni-tools/migrate.sh"}); err != nil {
		return nil, utils.Errorf(err, L("cannot run uyuni migration container"))
	}

	// now that everything is migrated, we need to fix SELinux permission
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

		pgsqlVersionUpgradeScriptName, err := adm_utils.GeneratePgsqlVersionUpgradeScript(scriptDir, oldPgsql, newPgsql)
		if err != nil {
			return utils.Errorf(err, L("cannot generate PostgreSQL database version upgrade script"))
		}

		err = podman.RunContainer(pgsqlVersionUpgradeContainer, preparedImage, utils.PgsqlRequiredVolumeMounts, extraArgs,
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
		"--network", podman.UyuniNetwork,
	}
	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"
	pgsqlFinalizeScriptName, err := adm_utils.GenerateFinalizePostgresScript(
		scriptDir, true, schemaUpdateRequired, migration, false,
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
	postUpgradeScriptName, err := adm_utils.GeneratePostUpgradeScript(scriptDir)
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
	systemd podman.Systemd,
	authFile string,
	registry string,
	db adm_utils.DBFlags,
	reportdb adm_utils.DBFlags,
	ssl adm_utils.InstallSSLFlags,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	cocoFlags adm_utils.CocoFlags,
	hubXmlrpcFlags adm_utils.HubXmlrpcFlags,
	salineFlags adm_utils.SalineFlags,
	pgsqlFlags types.PgsqlFlags,
	scc types.SCCCredentials,
	tz string,
) error {
	// Calling cloudguestregistryauth only makes sense if using the cloud provider registry.
	// This check assumes users won't use custom registries that are not the cloud provider one on a cloud image.
	if !strings.HasPrefix(registry, "registry.suse.com") {
		if err := CallCloudGuestRegistryAuth(); err != nil {
			return err
		}
	}

	// Prepare Uyuni network, migration container needs to run in the same network as resulting image
	err := podman.SetupNetwork(false)
	if err != nil {
		return utils.Errorf(err, L("cannot setup network"))
	}

	preparedServerImage, preparedPgsqlImage, err := podman.PrepareImages(authFile, image, pgsqlFlags)
	if err != nil {
		return utils.Errorf(err, L("cannot prepare images"))
	}

	inspectedValues, err := prepareHost(preparedServerImage, preparedPgsqlImage, image.PullPolicy, scc)
	if err != nil {
		return utils.Errorf(err, L("cannot prepare host"))
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

	oldPgVersion, _ := strconv.Atoi(inspectedValues.CommonInspectData.CurrentPgVersion)
	newPgVersion, _ := strconv.Atoi(inspectedValues.DBInspectData.ImagePgVersion)

	if inspectedValues.CommonInspectData.CurrentPgVersionNotMigrated != "" ||
		inspectedValues.DBHost == "localhost" ||
		inspectedValues.ReportDBHost == "localhost" {
		log.Info().Msgf(L("Configuring split PostgreSQL container. Image version: %[1]s, not migrated version: %[2]s"),
			newPgVersion, oldPgVersion)

		if err := configureSplitDBContainer(
			preparedServerImage, preparedPgsqlImage, systemd, db, reportdb, ssl, tz); err != nil {
			return utils.Errorf(err, L("cannot configure db container"))
		}
	}

	if newPgVersion > oldPgVersion {
		log.Info().Msgf(
			L("Previous PostgreSQL is %[1]s, instead new one is %[2]s. Performing a DB version upgrade…"),
			oldPgVersion, newPgVersion,
		)
		if err := RunPgsqlVersionUpgrade(
			authFile, registry, pgsqlFlags.Image, upgradeImage, strconv.Itoa(oldPgVersion),
			strconv.Itoa(newPgVersion),
		); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	} else if newPgVersion == oldPgVersion {
		log.Info().Msg(L("Upgrading without changing PostgreSQL version"))
	} else {
		return fmt.Errorf(
			L("trying to downgrade PostgreSQL from %[1]s to %[2]s"),
			oldPgVersion, newPgVersion,
		)
	}

	if err := pgsql.Upgrade(preparedPgsqlImage, systemd); err != nil {
		return err
	}

	schemaUpdateRequired :=
		oldPgVersion != newPgVersion
	if err := RunPgsqlFinalizeScript(preparedServerImage, schemaUpdateRequired, false); err != nil {
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

	if err := updateServerSystemdService(); err != nil {
		return err
	}
	log.Info().Msg(L("Waiting for the server to start…"))

	inspectedDB := adm_utils.DBFlags{
		Name:     inspectedValues.DBName,
		Port:     inspectedValues.DBPort,
		User:     inspectedValues.DBUser,
		Password: inspectedValues.DBPassword,
		Host:     db.Host,
	}

	err = coco.Upgrade(systemd, authFile, registry, cocoFlags, image, inspectedDB)

	if err != nil {
		return utils.Errorf(err, L("error upgrading confidential computing service."))
	}

	if err := hub.Upgrade(
		systemd, authFile, registry, image.PullPolicy, image.Tag, hubXmlrpcFlags,
	); err != nil {
		return err
	}

	if err := saline.Upgrade(systemd, authFile, registry, salineFlags, image,
		utils.GetLocalTimezone(), viper.GetStringSlice("podman.arg"),
	); err != nil {
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

	return cnx.WaitForServer()
}

// Migrate will migrate a server to the image given as attribute.
func Migrate(
	systemd podman.Systemd,
	authFile string,
	registry string,
	db adm_utils.DBFlags,
	reportdb adm_utils.DBFlags,
	ssl adm_utils.InstallSSLFlags,
	image types.ImageFlags,
	upgradeImage types.ImageFlags,
	cocoFlags adm_utils.CocoFlags,
	hubXmlrpcFlags adm_utils.HubXmlrpcFlags,
	salineFlags adm_utils.SalineFlags,
	pgsqlFlags types.PgsqlFlags,
	scc types.SCCCredentials,
	tz string,
	prepare bool,
	user string,
	debug bool,
	mirror string,
	podmanArgs podman.PodmanFlags,
	args []string,
) error {
	// Calling cloudguestregistryauth only makes sense if using the cloud provider registry.
	// This check assumes users won't use custom registries that are not the cloud provider one on a cloud image.
	if !strings.HasPrefix(registry, "registry.suse.com") {
		if err := CallCloudGuestRegistryAuth(); err != nil {
			return err
		}
	}

	sourceFqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}

	// Prepare Uyuni network, migration container needs to run in the same network as resulting image
	err = podman.SetupNetwork(false)
	if err != nil {
		return utils.Errorf(err, L("cannot setup network"))
	}
	// Find the SSH Socket and paths for the migration
	sshAuthSocket := GetSSHAuthSocket()
	sshConfigPath, sshKnownhostsPath := GetSSHPaths()

	preparedServerImage, preparedPgsqlImage, err := podman.PrepareImages(authFile, image, pgsqlFlags)
	if err != nil {
		return utils.Errorf(err, L("cannot prepare images"))
	}

	if err := stopService(systemd, podman.ServerService); err != nil {
		return err
	}
	if err := stopService(systemd, podman.DBService); err != nil {
		return err
	}

	_, err = RunMigration(
		preparedServerImage, sshAuthSocket, sshConfigPath, sshKnownhostsPath, sourceFqdn,
		user, prepare,
	)
	if err != nil {
		return utils.Errorf(err, L("cannot run migration script"))
	}
	if prepare {
		log.Info().Msg(L("Migration prepared. Run the 'migrate' command without '--prepare' to finish the migration."))
		return nil
	}

	inspectedValues, err := prepareHost(preparedServerImage, preparedPgsqlImage, image.PullPolicy, scc)
	if err != nil {
		return utils.Errorf(err, L("cannot prepare host"))
	}

	oldPgVersion, _ := strconv.Atoi("14")

	newPgVersion, _ := strconv.Atoi(inspectedValues.DBInspectData.ImagePgVersion)

	log.Info().Msgf(L("Configuring split PostgreSQL container. Image version: %[1]s, not migrated version: %[2]s"),
		newPgVersion, oldPgVersion)

	if err := upgradeDB(newPgVersion, oldPgVersion, upgradeImage, authFile, registry, pgsqlFlags.Image); err != nil {
		return err
	}

	if err := configureSplitDBContainer(
		preparedServerImage, preparedPgsqlImage, systemd, db, reportdb, ssl, tz); err != nil {
		return utils.Errorf(err, L("cannot configure db container"))
	}

	if err := pgsql.Upgrade(preparedPgsqlImage, systemd); err != nil {
		return err
	}

	schemaUpdateRequired :=
		oldPgVersion != newPgVersion
	if err := RunPgsqlFinalizeScript(preparedServerImage, schemaUpdateRequired, false); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL finalize script"))
	}

	if err := RunPostUpgradeScript(preparedServerImage); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	cnx := shared.NewConnection("podman", podman.ServerContainerName, "")
	if err := WaitForSystemStart(systemd, cnx, preparedServerImage, tz,
		debug, mirror, podmanArgs.Args); err != nil {
		return utils.Error(err, L("cannot wait for system start"))
	}

	inspectedDB := adm_utils.DBFlags{
		Name:     inspectedValues.DBName,
		Port:     inspectedValues.DBPort,
		User:     inspectedValues.DBUser,
		Password: inspectedValues.DBPassword,
		Host:     db.Host,
	}

	err = coco.Upgrade(systemd, authFile, registry, cocoFlags, image, inspectedDB)
	if err != nil {
		return utils.Errorf(err, L("error upgrading confidential computing service."))
	}

	if err := hub.Upgrade(
		systemd, authFile, registry, image.PullPolicy, image.Tag, hubXmlrpcFlags,
	); err != nil {
		return err
	}

	if err := saline.Upgrade(systemd, authFile, registry, salineFlags, image,
		utils.GetLocalTimezone(), viper.GetStringSlice("podman.arg"),
	); err != nil {
		return utils.Errorf(err, L("error upgrading saline service."))
	}

	return systemd.ReloadDaemon(false)
}

func stopService(systemd podman.Systemd, name string) error {
	if systemd.HasService(name) {
		if err := systemd.StopService(name); err != nil {
			return utils.Error(err, L("cannot stop service"))
		}
		defer func() {
			_ = systemd.StartService(name)
		}()
	}
	return nil
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

// RunPgsqlContainerMigration migrate to separate postgres container.
func RunPgsqlContainerMigration(serverImage string, dbHost string, reportDBHost string) error {
	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	data := templates.PgsqlMigrateScriptTemplateData{
		DBHost:       dbHost,
		ReportDBHost: reportDBHost,
	}

	scriptPath := filepath.Join(scriptDir, "pgmigrate.sh")
	if err = utils.WriteTemplateToFile(data, scriptPath, 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate postgresql migration script"))
	}

	podmanArgs := []string{
		"-v", scriptDir + ":" + scriptDir,
		"--security-opt", "label=disable",
	}
	err = podman.RunContainer("uyuni-db-migrate", serverImage, utils.ServerMigrationVolumeMounts, podmanArgs,
		[]string{scriptPath})

	return err
}

// RunPgsqlContainerMigration migrate to separate postgres container.
func RunConfigPgsl(pgsqlImage string) error {
	podmanArgs := []string{
		"--security-opt", "label=disable",
		"--entrypoint", "/docker-entrypoint-initdb.d/uyuni-postgres-config.sh",
	}
	if err := podman.RunContainer("uyuni-db-config", pgsqlImage, utils.ServerMigrationVolumeMounts,
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
	pullPolicy string,
	scc types.SCCCredentials,
) (*utils.ServerInspectData, error) {
	inspectedValues, err := podman.Inspect(preparedServerImage, preparedPgsqlImage, pullPolicy, scc)
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect podman values"))
	}

	runningImage := podman.GetServiceImage(podman.ServerService)
	var runningData *utils.ServerInspectData
	if runningImage != "" {
		runningData, err = podman.Inspect(preparedServerImage, preparedPgsqlImage, pullPolicy, scc)
		if err != nil {
			return inspectedValues, err
		}
	}

	return inspectedValues, adm_utils.SanityCheck(runningData, inspectedValues, preparedServerImage)
}

func upgradeDB(
	newPgVersion int,
	oldPgVersion int,
	upgradeImage types.ImageFlags,
	authFile string,
	registry string,
	dbImage types.ImageFlags,
) error {
	if newPgVersion > oldPgVersion {
		log.Info().Msgf(
			L("Previous PostgreSQL is %[1]s, instead new one is %[2]s. Performing a DB version upgrade…"),
			oldPgVersion, newPgVersion,
		)
		if err := RunPgsqlVersionUpgrade(
			authFile, registry, dbImage, upgradeImage, strconv.Itoa(oldPgVersion),
			strconv.Itoa(newPgVersion),
		); err != nil {
			return utils.Error(err, L("cannot run PostgreSQL version upgrade script"))
		}
	} else if newPgVersion == oldPgVersion {
		log.Info().Msg(L("Upgrading without changing PostgreSQL version"))
	} else {
		return fmt.Errorf(
			L("trying to downgrade PostgreSQL from %[1]s to %[2]s"),
			oldPgVersion, newPgVersion,
		)
	}
	return nil
}

func configureSplitDBContainer(
	serverImage string,
	pgsqlImage string,
	systemd podman.Systemd,
	db adm_utils.DBFlags,
	reportdb adm_utils.DBFlags,
	ssl adm_utils.InstallSSLFlags,
	tz string,
) error {
	if err := RunPgsqlContainerMigration(serverImage, "db", "reportdb"); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
	}
	fqdn, err := utils.GetFqdn([]string{})
	if err != nil {
		return err
	}

	if err = PrepareSSLCertificates(serverImage, &ssl, tz, fqdn); err != nil {
		return err
	}

	// Create all the database credentials secrets
	if err := podman.CreateCredentialsSecrets(
		podman.DBUserSecret, db.User,
		podman.DBPassSecret, db.Password,
	); err != nil {
		return err
	}

	if err := podman.CreateCredentialsSecrets(
		podman.ReportDBUserSecret, reportdb.User,
		podman.ReportDBPassSecret, reportdb.Password,
	); err != nil {
		return err
	}

	if db.IsLocal() {
		// The admin password is not needed for external databases
		if err := podman.CreateCredentialsSecrets(
			podman.DBAdminUserSecret, db.Admin.User,
			podman.DBAdminPassSecret, db.Admin.Password,
		); err != nil {
			return err
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
