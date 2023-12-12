// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToKubernetes(globalFlags *types.GlobalFlags, flags *kubernetesMigrateFlags, cmd *cobra.Command, args []string) {
	cnx := utils.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)
	fqdn := args[0]

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := shared.GetSshPaths()

	// Prepare the migration script and folder
	scriptDir := shared.GenerateMigrationScript(fqdn, true)
	defer os.RemoveAll(scriptDir)

	// Install Uyuni with generated CA cert: an empty struct means no 3rd party cert
	var sslFlags adm_utils.SslCertFlags

	// We don't need the SSL certs at this point of the migration
	clusterInfos := shared_kubernetes.CheckCluster()

	kubeconfig := clusterInfos.GetKubeconfig()
	//TODO: check if we need to handle SELinux policies, as we do in podman

	kubernetes.Deploy(cnx, &flags.Image, &flags.Helm, &sslFlags, &clusterInfos, fqdn, false,
		"--set", "migration.ssh.agentSocket="+sshAuthSocket,
		"--set", "migration.ssh.configPath="+sshConfigPath,
		"--set", "migration.ssh.knownHostsPath="+sshKnownhostsPath,
		"--set", "migration.dataPath="+scriptDir)

	// Run the actual migration
	runMigration(cnx, flags, scriptDir)

	tz, oldPgVersion, newPgVersion := shared.ReadContainerData(scriptDir)

	helmArgs := []string{
		"--reset-values",
		"--set", "timezone=" + tz,
	}

	if oldPgVersion != newPgVersion {
		var migrationImage types.ImageFlags
		migrationImage.Name = flags.MigrationImage.Name
		if migrationImage.Name == "" {
			migrationImage.Name = fmt.Sprintf("%s-migration-%s-%s", flags.Image.Name, oldPgVersion, newPgVersion)
		}
		migrationImage.Tag = flags.MigrationImage.Tag
		log.Info().Msgf("Using migration image %s:%s", migrationImage.Name, migrationImage.Tag)
		shared.GeneratePgMigrationScript(scriptDir, oldPgVersion, newPgVersion, false)

		kubernetes.UyuniUpgrade(&migrationImage, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
		runMigration(cnx, flags, scriptDir)
	}

	shared.GenerateFinalizePostgresMigrationScript(scriptDir, true, oldPgVersion != newPgVersion, true, true, false)
	kubernetes.UyuniUpgrade(&flags.Image, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
	runMigration(cnx, flags, scriptDir)

	helmArgs = append(helmArgs, setupSsl(flags, kubeconfig, scriptDir)...)

	// As we upgrade the helm instance without the migration parameters the SSL certificate will be used
	kubernetes.UyuniUpgrade(&flags.Image, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
}

func runMigration(cnx *utils.Connection, flags *kubernetesMigrateFlags, tmpPath string) {
	log.Info().Msg("Migrating server")
	err := adm_utils.ExecCommand(zerolog.InfoLevel, cnx, "/var/lib/uyuni-tools/migrate.sh")
	if err != nil {
		log.Fatal().Err(err).Msg("error running the migration script")
	}
}

// updateIssuer replaces the temporary SSL certificate issuer with the source server CA.
// Return additional helm args to use the SSL certificates
func setupSsl(flags *kubernetesMigrateFlags, kubeconfig string, scriptDir string) []string {
	caCert := path.Join(scriptDir, "RHN-ORG-TRUSTED-SSL-CERT")
	caKey := path.Join(scriptDir, "RHN-ORG-PRIVATE-SSL-KEY")

	if utils.FileExists(caCert) && utils.FileExists(caKey) {
		key := base64.StdEncoding.EncodeToString(ssl.GetRsaKey(caKey, flags.Ssl.Password))

		// Strip down the certificate text part
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "openssl", "x509", "-in", caCert)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to strip text part of CA certificate")
		}
		cert := base64.StdEncoding.EncodeToString(out)
		ca := ssl.SslPair{Cert: cert, Key: key}

		// An empty struct means no third party certificate
		sslFlags := adm_utils.SslCertFlags{}
		return kubernetes.DeployCertificate(&flags.Helm, &sslFlags, cert, &ca, kubeconfig, "", flags.Image.PullPolicy)
	} else {
		// Handle third party certificates and CA
		sslFlags := adm_utils.SslCertFlags{
			Ca: ssl.CaChain{Root: caCert},
			Server: ssl.SslPair{
				Key:  path.Join(scriptDir, "spacewalk.key"),
				Cert: path.Join(scriptDir, "spacewalk.crt"),
			},
		}
		kubernetes.DeployExistingCertificate(&flags.Helm, &sslFlags, kubeconfig)
	}
	return []string{}
}
