// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

func migrateToKubernetes(globalFlags *types.GlobalFlags, flags *kubernetesMigrateFlags, cmd *cobra.Command, args []string) {
	cnx := utils.NewConnection("kubectl")
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
	clusterInfos := kubernetes.CheckCluster()

	kubernetes.Deploy(cnx, &flags.Image, &flags.Helm, &sslFlags, &clusterInfos, fqdn, false,
		"--set", "migration.ssh.agentSocket="+sshAuthSocket,
		"--set", "migration.ssh.configPath="+sshConfigPath,
		"--set", "migration.ssh.knownHostsPath="+sshKnownhostsPath,
		"--set", "migration.dataPath="+scriptDir)

	// Run the actual migration
	runMigration(cnx, flags, scriptDir)

	tz := shared.ReadTimezone(scriptDir)

	helmArgs := []string{
		"--reset-values",
		"--set", "timezone=" + tz,
	}

	kubeconfig := clusterInfos.GetKubeconfig()
	helmArgs = append(helmArgs, setupSsl(flags, kubeconfig, scriptDir)...)

	// As we upgrade the helm instance without the migration parameters the SSL certificate will be used
	kubernetes.UyuniUpgrade(&flags.Image, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
}

func runMigration(cnx *utils.Connection, flags *kubernetesMigrateFlags, tmpPath string) {
	log.Info().Msg("Migrating server")
	err := adm_utils.ExecCommand(zerolog.DebugLevel, cnx, "/var/lib/uyuni-tools/migrate.sh")
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
