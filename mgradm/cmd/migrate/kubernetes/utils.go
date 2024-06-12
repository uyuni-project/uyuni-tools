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
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	migration_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/migrate/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
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
			return fmt.Errorf(L("install %s before running this command"), binary)
		}
	}
	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)

	serverImage, err := utils.ComputeImage(flags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	fqdn := args[0]

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := migration_shared.GetSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := migration_shared.GetSshPaths()

	// Prepare the migration script and folder
	scriptDir, err := adm_utils.GenerateMigrationScript(fqdn, flags.User, true)
	if err != nil {
		return utils.Errorf(err, L("failed to generate migration script"))
	}

	defer os.RemoveAll(scriptDir)

	// We don't need the SSL certs at this point of the migration
	clusterInfos, err := shared_kubernetes.CheckCluster()
	if err != nil {
		return err
	}
	kubeconfig := clusterInfos.GetKubeconfig()
	//TODO: check if we need to handle SELinux policies, as we do in podman

	// Install Uyuni with generated CA cert: an empty struct means no 3rd party cert
	var sslFlags adm_utils.SslCertFlags

	// Deploy for running migration command
	if err := kubernetes.Deploy(cnx, &flags.Image, &flags.Helm, &sslFlags, clusterInfos, fqdn, false,
		"--set", "migration.ssh.agentSocket="+sshAuthSocket,
		"--set", "migration.ssh.configPath="+sshConfigPath,
		"--set", "migration.ssh.knownHostsPath="+sshKnownhostsPath,
		"--set", "migration.dataPath="+scriptDir); err != nil {
		return utils.Errorf(err, L("cannot run deploy"))
	}

	//this is needed because folder with script needs to be mounted
	//check the node before scaling down
	nodeName, err := shared_kubernetes.GetNode("uyuni")
	if err != nil {
		return utils.Errorf(err, L("cannot find node running uyuni"))
	}
	// Run the actual migration
	if err := adm_utils.RunMigration(cnx, scriptDir, "migrate.sh"); err != nil {
		return utils.Errorf(err, L("cannot run migration"))
	}

	tz, oldPgVersion, newPgVersion, err := adm_utils.ReadContainerData(scriptDir)
	if err != nil {
		return utils.Errorf(err, L("cannot read data from container"))
	}

	// After each command we want to scale to 0
	err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerApp, 0)
	if err != nil {
		return utils.Errorf(err, L("cannot set replicas to 0"))
	}

	defer func() {
		// if something is running, we don't need to set replicas to 1
		if _, err = shared_kubernetes.GetNode("uyuni"); err != nil {
			err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerApp, 1)
		}
	}()

	setupSslArray, err := setupSsl(&flags.Helm, kubeconfig, scriptDir, flags.Ssl.Password, flags.Image.PullPolicy)
	if err != nil {
		return utils.Errorf(err, L("cannot setup SSL"))
	}

	helmArgs := []string{
		"--reset-values",
		"--set", "timezone=" + tz,
	}
	if flags.Mirror != "" {
		log.Warn().Msgf(L("The mirror data will not be migrated, ensure it is available at %s"), flags.Mirror)
		// TODO Handle claims for multi-node clusters
		helmArgs = append(helmArgs, "--set", "mirror.hostPath="+flags.Mirror)
	}
	helmArgs = append(helmArgs, setupSslArray...)

	// Run uyuni upgrade using the new ssl certificate
	err = kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
	if err != nil {
		return utils.Errorf(err, L("cannot upgrade helm chart to image %s using new SSL certificate"), serverImage)
	}

	if err := shared_kubernetes.WaitForDeployment(flags.Helm.Uyuni.Namespace, "uyuni", "uyuni"); err != nil {
		return utils.Errorf(err, L("cannot wait for deployment of %s"), serverImage)
	}

	err = shared_kubernetes.ReplicasTo(shared_kubernetes.ServerApp, 0)
	if err != nil {
		return utils.Errorf(err, L("cannot set replicas to 0"))
	}

	if oldPgVersion != newPgVersion {
		if err := kubernetes.RunPgsqlVersionUpgrade(flags.Image, flags.DbUpgradeImage, nodeName, oldPgVersion, newPgVersion); err != nil {
			return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
		}
	}

	schemaUpdateRequired := oldPgVersion != newPgVersion
	if err := kubernetes.RunPgsqlFinalizeScript(serverImage, flags.Image.PullPolicy, nodeName, schemaUpdateRequired); err != nil {
		return utils.Errorf(err, L("cannot run PostgreSQL version upgrade script"))
	}

	if err := kubernetes.RunPostUpgradeScript(serverImage, flags.Image.PullPolicy, nodeName); err != nil {
		return utils.Errorf(err, L("cannot run post upgrade script"))
	}

	err = kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
	if err != nil {
		return utils.Errorf(err, L("cannot upgrade to image %s"), serverImage)
	}

	if err := shared_kubernetes.WaitForDeployment(flags.Helm.Uyuni.Namespace, "uyuni", "uyuni"); err != nil {
		return err
	}

	if err := cnx.CopyCaCertificate(fqdn); err != nil {
		return utils.Errorf(err, L("failed to add SSL CA certificate to host trusted certificates"))
	}
	return nil
}

// updateIssuer replaces the temporary SSL certificate issuer with the source server CA.
// Return additional helm args to use the SSL certificates.
func setupSsl(helm *adm_utils.HelmFlags, kubeconfig string, scriptDir string, password string, pullPolicy string) ([]string, error) {
	caCert := path.Join(scriptDir, "RHN-ORG-TRUSTED-SSL-CERT")
	caKey := path.Join(scriptDir, "RHN-ORG-PRIVATE-SSL-KEY")

	if utils.FileExists(caCert) && utils.FileExists(caKey) {
		key := base64.StdEncoding.EncodeToString(ssl.GetRsaKey(caKey, password))

		// Strip down the certificate text part
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "openssl", "x509", "-in", caCert)
		if err != nil {
			return []string{}, utils.Errorf(err, L("failed to strip text part from CA certificate"))
		}
		cert := base64.StdEncoding.EncodeToString(out)
		ca := ssl.SslPair{Cert: cert, Key: key}

		// An empty struct means no third party certificate
		sslFlags := adm_utils.SslCertFlags{}
		ret, err := kubernetes.DeployCertificate(helm, &sslFlags, cert, &ca, kubeconfig, "", pullPolicy)
		if err != nil {
			return []string{}, utils.Errorf(err, L("cannot deploy certificate"))
		}
		return ret, nil
	} else {
		// Handle third party certificates and CA
		sslFlags := adm_utils.SslCertFlags{
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
