package migrate

import (
	"encoding/base64"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyuniadm/shared/utils"
)

func migrateToKubernetes(globalFlags *types.GlobalFlags, flags *MigrateFlags, cmd *cobra.Command, args []string) {
	fqdn := args[0]

	// Find the SSH Socket and paths for the migration
	sshAuthSocket := getSshAuthSocket()
	sshConfigPath, sshKnownhostsPath := getSshPaths()

	// Prepare the migration script and folder
	scriptDir := generateMigrationScript(fqdn, true)
	defer os.RemoveAll(scriptDir)

	// Install Uyuni with generated CA cert
	var sslFlags cmd_utils.SslCertFlags
	sslFlags.UseExisting = false

	// We don't need the SSL certs at this point of the migration
	clusterInfos := kubernetes.CheckCluster()

	kubernetes.Deploy(globalFlags, &flags.Image, &flags.Helm, &sslFlags, &clusterInfos, fqdn,
		"--set", "migration.ssh.agentSocket="+sshAuthSocket,
		"--set", "migration.ssh.configPath="+sshConfigPath,
		"--set", "migration.ssh.knownHostsPath="+sshKnownhostsPath,
		"--set", "migration.dataPath="+scriptDir)

	// Run the actual migration
	runMigration(globalFlags, flags, scriptDir)

	tz := readTimezone(scriptDir)

	helmArgs := []string{
		"--reset-values",
		"--set", "timezone=" + tz,
	}

	// TODO Update uyuni-ca secret with the source CA cert and key
	kubeconfig := clusterInfos.GetKubeconfig()
	helmArgs = append(helmArgs, setupSsl(globalFlags, flags, kubeconfig, scriptDir)...)

	// Update the installation in non-migration mode

	// As we upgrade the helm instance without the migration parameters the SSL certificate will be used
	kubernetes.UyuniUpgrade(globalFlags, &flags.Image, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress, helmArgs...)
}

func runMigration(globalFlags *types.GlobalFlags, flags *MigrateFlags, tmpPath string) {
	log.Println("Migrating server")
	utils.Exec(globalFlags, "", false, false, []string{}, "/var/lib/uyuni-tools/migrate.sh")
}

// updateIssuer replaces the temporary SSL certificate issuer with the source server CA.
// Return additional helm args to use the SSL certificates
func setupSsl(globalFlags *types.GlobalFlags, flags *MigrateFlags, kubeconfig string, scriptDir string) []string {
	caCert := path.Join(scriptDir, "RHN-ORG-TRUSTED-SSL-CERT")
	caKey := path.Join(scriptDir, "RHN-ORG-PRIVATE-SSL-KEY")

	if utils.FileExists(caCert) && utils.FileExists(caKey) {
		// Convert the key file to RSA format for kubectl to handle it
		cmd := exec.Command("openssl", "rsa", "-in", caKey, "-passin", "env:pass")
		cmd.Env = append(cmd.Env, "pass=spacewalk") // TODO Parametrize!
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("Failed to convert CA private key to RSA: %s\n", err)
		}
		key := base64.StdEncoding.EncodeToString(out)

		// Strip down the certificate text part
		out, err = exec.Command("openssl", "x509", "-in", caCert).Output()
		if err != nil {
			log.Fatalf("Failed to strip text part of CA certificate: %s\n", err)
		}
		cert := base64.StdEncoding.EncodeToString(out)
		ca := kubernetes.TlsCert{RootCa: cert, Certificate: cert, Key: key}

		sslFlags := cmd_utils.SslCertFlags{UseExisting: false}
		return kubernetes.DeployCertificate(globalFlags, &flags.Helm, &sslFlags, &ca, kubeconfig, "")
	} else {
		// TODO Handle third party certificates and CA
	}
	return []string{}
}
