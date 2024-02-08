// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	migration_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func migrateToKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetesMigrateFlags,
	cmd *cobra.Command,
	args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf("install %s before running this command: %s", binary, err)
		}
	}
	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)
	fqdn := args[0]

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	// Prepare the migration script and folder
	scriptDir, err := adm_utils.GenerateMigrationScript(fqdn, true)
	if err != nil {
		return fmt.Errorf("Failed to generater migration script: %s", err)
	}

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
	adm_utils.RunMigration(cnx, scriptDir, "migrate.sh")

	tz, oldPgVersion, newPgVersion, err := adm_utils.ReadContainerData(scriptDir)
	if err != nil {
		return fmt.Errorf("cannot read data from container: %s", err)
	}

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

		scriptName, err := adm_utils.GeneratePgMigrationScript(scriptDir, oldPgVersion, newPgVersion, false)
		if err != nil {
			return fmt.Errorf("Cannot generate postgresql database migration script: %s", err)
		}

		migrationImageUrl, err := utils.ComputeImage(migrationImage.Name, migrationImage.Tag)
		if err != nil {
			return fmt.Errorf("Failed to compute image URL")
		}

		kubernetes.UyuniUpgrade(migrationImageUrl, migrationImage.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
		adm_utils.RunMigration(cnx, scriptDir, scriptName)
	}

	scriptName, err := adm_utils.GenerateFinalizePostgresMigrationScript(scriptDir, true, oldPgVersion != newPgVersion, true, true, false)
	if err != nil {
		return fmt.Errorf("Cannot generate postgresql migration finalization script: %s", err)
	}

	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("Failed to compute image URL")
	}

	kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
	adm_utils.RunMigration(cnx, scriptDir, scriptName)

	setupSslArray, err := setupSsl(&flags.Helm, kubeconfig, scriptDir, flags.Ssl.Password, flags.Image.PullPolicy)
	if err != nil {
		return fmt.Errorf("Cannot setup SSL: %s", err)
	}
	helmArgs = append(helmArgs, setupSslArray...)

	// As we upgrade the helm instance without the migration parameters the SSL certificate will be used
	kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
	return nil
}

// updateIssuer replaces the temporary SSL certificate issuer with the source server CA.
// Return additional helm args to use the SSL certificates
func setupSsl(helm *cmd_utils.HelmFlags, kubeconfig string, scriptDir string, password string, pullPolicy string) ([]string, error) {
	caCert := path.Join(scriptDir, "RHN-ORG-TRUSTED-SSL-CERT")
	caKey := path.Join(scriptDir, "RHN-ORG-PRIVATE-SSL-KEY")

	if utils.FileExists(caCert) && utils.FileExists(caKey) {
		key := base64.StdEncoding.EncodeToString(ssl.GetRsaKey(caKey, password))

		// Strip down the certificate text part
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "openssl", "x509", "-in", caCert)
		if err != nil {
			return []string{}, fmt.Errorf("Failed to strip text part of CA certificate %s", err)
		}
		cert := base64.StdEncoding.EncodeToString(out)
		ca := ssl.SslPair{Cert: cert, Key: key}

		// An empty struct means no third party certificate
		sslFlags := cmd_utils.SslCertFlags{}
		return kubernetes.DeployCertificate(helm, &sslFlags, cert, &ca, kubeconfig, "", pullPolicy), nil
	} else {
		// Handle third party certificates and CA
		sslFlags := cmd_utils.SslCertFlags{
			Ca: ssl.CaChain{Root: caCert},
			Server: ssl.SslPair{
				Key:  path.Join(scriptDir, "spacewalk.key"),
				Cert: path.Join(scriptDir, "spacewalk.crt"),
			},
		}
		kubernetes.DeployExistingCertificate(helm, &sslFlags, kubeconfig)
	}
	return []string{}, nil
}
