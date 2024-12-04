// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/saline"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func waitForSystemStart(
	systemd shared_podman.Systemd,
	cnx *shared.Connection,
	image string,
	flags *podmanInstallFlags,
) error {
	err := podman.GenerateSystemdService(
		systemd, flags.Installation.TZ, image, flags.Installation.Debug.Java, flags.Mirror, flags.Podman.Args,
	)
	if err != nil {
		return err
	}

	log.Info().Msg(L("Waiting for the server to start…"))
	if err := systemd.EnableService(shared_podman.ServerService); err != nil {
		return utils.Errorf(err, L("cannot enable service"))
	}

	return cnx.WaitForServer()
}

var systemd shared_podman.Systemd = shared_podman.SystemdImpl{}

func installForPodman(
	_ *types.GlobalFlags,
	flags *podmanInstallFlags,
	cmd *cobra.Command,
	args []string,
) error {
	hostData, err := shared_podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := shared_podman.PodmanLogin(hostData, flags.Installation.SCC)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	if hostData.HasUyuniServer {
		return errors.New(
			L("Server is already initialized! Uninstall before attempting new installation or use upgrade command"),
		)
	}

	flags.Installation.CheckParameters(cmd, "podman")
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	fqdn, err := utils.GetFqdn(args)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("Setting up the server with the FQDN '%s'"), fqdn)

	image, err := utils.ComputeImage(flags.Image.Registry, utils.DefaultTag, flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := shared_podman.PrepareImage(authFile, image, flags.Image.PullPolicy, true)
	if err != nil {
		return err
	}

	if err := shared_podman.SetupNetwork(false); err != nil {
		return utils.Errorf(err, L("cannot setup network"))
	}

	log.Info().Msg(L("Run setup command in the container"))

	if err := runSetup(preparedImage, &flags.ServerFlags, fqdn); err != nil {
		return err
	}

	cnx := shared.NewConnection("podman", shared_podman.ServerContainerName, "")
	if err := waitForSystemStart(systemd, cnx, preparedImage, flags); err != nil {
		return utils.Errorf(err, L("cannot wait for system start"))
	}

	if err := cnx.CopyCaCertificate(fqdn); err != nil {
		return utils.Errorf(err, L("failed to add SSL CA certificate to host trusted certificates"))
	}

	if path, err := exec.LookPath("uyuni-payg-extract-data"); err == nil {
		// the binary is installed
		err = utils.RunCmdStdMapping(zerolog.DebugLevel, path)
		if err != nil {
			return utils.Errorf(err, L("failed to extract payg data"))
		}
	}

	if flags.Coco.Replicas > 0 {
		// This may need to be moved up later once more containers require DB access
		if err := shared_podman.CreateDBSecrets(flags.Installation.DB.User, flags.Installation.DB.Password); err != nil {
			return err
		}
		if err := coco.SetupCocoContainer(
			systemd, authFile, flags.Image.Registry, flags.Coco, flags.Image,
			flags.Installation.DB.Name, flags.Installation.DB.Port,
		); err != nil {
			return err
		}
	}

	if flags.HubXmlrpc.Replicas > 0 {
		if err := hub.SetupHubXmlrpc(
			systemd, authFile, flags.Image.Registry, flags.Image.PullPolicy, flags.Image.Tag, flags.HubXmlrpc,
		); err != nil {
			return err
		}
	}

	if flags.Saline.Replicas > 0 {
		if err := saline.SetupSalineContainer(
			systemd, authFile, flags.Image.Registry, flags.Saline, flags.Image,
			flags.Installation.TZ, flags.Podman.Args,
		); err != nil {
			return err
		}
	}

	if flags.Installation.SSL.UseExisting() {
		if err := podman.UpdateSSLCertificate(
			cnx, &flags.Installation.SSL.Ca, &flags.Installation.SSL.Server,
		); err != nil {
			return utils.Errorf(err, L("cannot update SSL certificate"))
		}
	}

	if err := shared_podman.EnablePodmanSocket(); err != nil {
		return utils.Errorf(err, L("cannot enable podman socket"))
	}
	return nil
}

// runSetup execute the setup.
func runSetup(image string, flags *adm_utils.ServerFlags, fqdn string) error {
	localHostValues := []string{
		"localhost",
		"127.0.0.1",
		"::1",
		fqdn,
	}

	localDB := utils.Contains(localHostValues, flags.Installation.DB.Host)

	dbHost := flags.Installation.DB.Host
	reportdbHost := flags.Installation.ReportDB.Host

	if localDB {
		dbHost = "localhost"
		if reportdbHost == "" {
			reportdbHost = "localhost"
		}
	}

	caPassword := flags.Installation.SSL.Password
	if flags.Installation.SSL.UseExisting() {
		// We need to have a password for the generated CA, even though it will be thrown away after install
		caPassword = "dummy"
	}

	// TODO Share the env variables preparation with Kubernetes?
	env := map[string]string{
		"UYUNI_FQDN":            fqdn,
		"MANAGER_USER":          flags.Installation.DB.User,
		"MANAGER_PASS":          flags.Installation.DB.Password,
		"ADMIN_USER":            flags.Installation.Admin.Login,
		"ADMIN_PASS":            flags.Installation.Admin.Password,
		"MANAGER_ADMIN_EMAIL":   flags.Installation.Email,
		"MANAGER_MAIL_FROM":     flags.Installation.EmailFrom,
		"MANAGER_ENABLE_TFTP":   boolToString(flags.Installation.Tftp),
		"LOCAL_DB":              boolToString(localDB),
		"MANAGER_DB_NAME":       flags.Installation.DB.Name,
		"MANAGER_DB_HOST":       dbHost,
		"MANAGER_DB_PORT":       strconv.Itoa(flags.Installation.DB.Port),
		"MANAGER_DB_PROTOCOL":   flags.Installation.DB.Protocol,
		"REPORT_DB_NAME":        flags.Installation.ReportDB.Name,
		"REPORT_DB_HOST":        reportdbHost,
		"REPORT_DB_PORT":        strconv.Itoa(flags.Installation.ReportDB.Port),
		"REPORT_DB_USER":        flags.Installation.ReportDB.User,
		"REPORT_DB_PASS":        flags.Installation.ReportDB.Password,
		"EXTERNALDB_ADMIN_USER": flags.Installation.DB.Admin.User,
		"EXTERNALDB_ADMIN_PASS": flags.Installation.DB.Admin.Password,
		"EXTERNALDB_PROVIDER":   flags.Installation.DB.Provider,
		"ISS_PARENT":            flags.Installation.IssParent,
		"ACTIVATE_SLP":          "N", // Deprecated, will be removed soon
		"SCC_USER":              flags.Installation.SCC.User,
		"SCC_PASS":              flags.Installation.SCC.Password,
		"CERT_O":                flags.Installation.SSL.Org,
		"CERT_OU":               flags.Installation.SSL.OU,
		"CERT_CITY":             flags.Installation.SSL.City,
		"CERT_STATE":            flags.Installation.SSL.State,
		"CERT_COUNTRY":          flags.Installation.SSL.Country,
		"CERT_EMAIL":            flags.Installation.SSL.Email,
		"CERT_CNAMES":           strings.Join(append([]string{fqdn}, flags.Installation.SSL.Cnames...), ","),
		"CERT_PASS":             caPassword,
	}

	if flags.Mirror != "" {
		env["MIRROR_PATH"] = "/mirror"
	}

	envNames := []string{}
	envValues := []string{}
	for key, value := range env {
		envNames = append(envNames, "-e", key)
		envValues = append(envValues, fmt.Sprintf("%s=%s", key, value))
	}

	command := []string{
		"run",
		"--rm",
		"--shm-size=0",
		"--shm-size-systemd=0",
		"--name", "uyuni-setup",
		"--network", shared_podman.UyuniNetwork,
		"-e", "TZ=" + flags.Installation.TZ,
	}
	for _, volume := range utils.ServerVolumeMounts {
		command = append(command, "-v", fmt.Sprintf("%s:%s:z", volume.Name, volume.MountPath))
	}
	command = append(command, envNames...)
	command = append(command, image)

	script, err := generateSetupScript(&flags.Installation)
	if err != nil {
		return err
	}
	command = append(command, "/usr/bin/sh", "-c", script)

	if _, err := newRunner("podman", command...).Env(envValues).StdMapping().Exec(); err != nil {
		return utils.Errorf(err, L("server setup failed"))
	}

	log.Info().Msgf(L("Server set up, login on https://%[1]s with %[2]s user"), fqdn, flags.Installation.Admin.Login)
	return nil
}

var newRunner = utils.NewRunner

// generateSetupScript creates a temporary folder with the setup script to execute in the container.
// The script exports all the needed environment variables and calls uyuni's mgr-setup.
func generateSetupScript(flags *adm_utils.InstallationFlags) (string, error) {
	// TODO Share with kubernetes implementation
	template := templates.MgrSetupScriptTemplateData{
		DebugJava:      flags.Debug.Java,
		OrgName:        flags.Organization,
		AdminLogin:     "$ADMIN_USER",
		AdminPassword:  "$ADMIN_PASS",
		AdminFirstName: flags.Admin.FirstName,
		AdminLastName:  flags.Admin.LastName,
		AdminEmail:     flags.Admin.Email,
		NoSSL:          false,
	}

	// Prepare the script
	scriptBuilder := new(strings.Builder)
	if err := template.Render(scriptBuilder); err != nil {
		return "", utils.Errorf(err, L("failed to render setup script"))
	}
	return scriptBuilder.String(), nil
}

func boolToString(value bool) string {
	if value {
		return "Y"
	}
	return "N"
}
