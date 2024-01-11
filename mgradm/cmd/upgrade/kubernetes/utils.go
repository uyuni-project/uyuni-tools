//go:build !nok8s

package kubernetes

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect"
	upgrade_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/upgrade/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/kubernetes"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func upgradeKubernetes(
	globalFlags *types.GlobalFlags,
	flags *kubernetesUpgradeFlags,
	cmd *cobra.Command,
	args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf("install %s before running this command", binary)
		}
	}
	cnx := shared.NewConnection("kubectl", "", shared_kubernetes.ServerFilter)
	if (len(args[0])) <= 0 {
		return fmt.Errorf("FQDN must be provided as argument")
	}
	fqdn := args[0]

	serverImage, err := utils.ComputeImage(flags.Image.Name, flags.Image.Tag)
	if err != nil {
		return fmt.Errorf("Failed to compute image URL")
	}

	inspectedValues, err := inspect.InspectKubernetes(serverImage, flags.Image.PullPolicy)

	upgrade_shared.SanityCheck(cnx, inspectedValues, serverImage)

	clusterInfos := shared_kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)

	if err != nil {
		return fmt.Errorf("Failed to create temporary directory")
	}

	shared_kubernetes.ReplicasTo(shared_kubernetes.ServerFilter, 0)
	defer shared_kubernetes.ReplicasTo(shared_kubernetes.ServerFilter, 1)

	nodeName, err := shared_kubernetes.GetNode("uyuni")

	if err != nil {
		return fmt.Errorf("Cannot find node for app uyuni %s", err)
	}

	pgsqlMigrationArgs := []string{}

	if inspectedValues["image_pg_version"] > inspectedValues["current_pg_version"] {
		log.Info().Msgf("Previous postgresql is %s, instead new one is %s. Performing a DB migration...", inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])

		migrationContainer := "uyuni-upgrade-pgsql"

		var migrationImage types.ImageFlags
		migrationImage.Name = flags.MigrationImage.Name
		if migrationImage.Name == "" {
			migrationImage.Name = fmt.Sprintf("%s-migration-%s-%s", flags.Image.Name, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"])
		}
		migrationImage.Tag = flags.MigrationImage.Tag
		migrationImageUrl, err := utils.ComputeImage(migrationImage.Name, flags.Image.Tag)
		if err != nil {
			return fmt.Errorf("Failed to compute image URL %s", err)
		}

		log.Info().Msgf("Using migration image %s", migrationImageUrl)

		//FIXME create a type for it and pass it as Go Template
		pgsqlMigrationArgs = []string{
			"--override-type=strategic",
			"--overrides",
			fmt.Sprintf(`{"apiVersion":"v1","spec":{"nodeName":"%s","restartPolicy":"Never","containers":[{"name":%s,"volumeMounts":[{"mountPath":"/etc/pki/tls","name":"etc-tls"},{"mountPath":"/var/lib/pgsql","name":"var-pgsql"},{"mountPath":"/var/lib/uyuni-tools","name":"var-lib-uyuni-tools"},{"mountPath":"/etc/rhn","name":"etc-rhn"},{"mountPath":"/etc/pki/spacewalk-tls","name":"tls-key"}]}],"volumes":[{"name":"etc-tls","persistentVolumeClaim":{"claimName":"etc-tls"}},{"name":"var-pgsql","persistentVolumeClaim":{"claimName":"var-pgsql"}},{"name":"var-lib-uyuni-tools","hostPath":{"path":%s,"type":"Directory"}},{"name":"etc-rhn","persistentVolumeClaim":{"claimName":"etc-rhn"}},{"name":"tls-key","secret":{"secretName":"uyuni-cert","items":[{"key":"tls.crt","path":"spacewalk.crt"},{"key":"tls.key","path":"spacewalk.key"}]}}]}}`,
				nodeName,
				strconv.Quote(migrationContainer),
				strconv.Quote(scriptDir),
			),
		}
		scriptName, err := adm_utils.GeneratePgMigrationScript(scriptDir, inspectedValues["current_pg_version"], inspectedValues["image_pg_version"], true)
		if err != nil {
			return fmt.Errorf("Cannot generate pg migration script: %s", err)
		}

		err = shared_kubernetes.RunPod(migrationContainer, migrationImageUrl, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+scriptName, pgsqlMigrationArgs...)
		if err != nil {
			return fmt.Errorf("Error running container %s: %s", migrationContainer, err)
		}
	}

	//FIXME finalize pgsql should be done automatically

	pgsqlFinalizeContainer := "uyuni-finalize-pgsql"

	pgsqlFinalizeArgs := []string{
		"--override-type=strategic",
		"--overrides",
		fmt.Sprintf(`{"apiVersion":"v1","spec":{"nodeName":"%s","restartPolicy":"Never","containers":[{"name":%s,"volumeMounts":[{"mountPath":"/etc/pki/tls","name":"etc-tls"},{"mountPath":"/var/lib/pgsql","name":"var-pgsql"},{"mountPath":"/var/lib/uyuni-tools","name":"var-lib-uyuni-tools"},{"mountPath":"/etc/rhn","name":"etc-rhn"},{"mountPath":"/etc/pki/spacewalk-tls","name":"tls-key"}]}],"volumes":[{"name":"etc-tls","persistentVolumeClaim":{"claimName":"etc-tls"}},{"name":"var-pgsql","persistentVolumeClaim":{"claimName":"var-pgsql"}},{"name":"var-lib-uyuni-tools","hostPath":{"path":%s,"type":"Directory"}},{"name":"etc-rhn","persistentVolumeClaim":{"claimName":"etc-rhn"}},{"name":"tls-key","secret":{"secretName":"uyuni-cert","items":[{"key":"tls.crt","path":"spacewalk.crt"},{"key":"tls.key","path":"spacewalk.key"}]}}]}}`,
			nodeName,
			strconv.Quote(pgsqlFinalizeContainer),
			strconv.Quote(scriptDir),
		),
	}

	scriptName, err := adm_utils.GenerateFinalizePostgresMigrationScript(scriptDir, true, inspectedValues["current_pg_version"] != inspectedValues["image_pg_version"], true, true, true)

	if err != nil {
		return fmt.Errorf("Cannot generate pg migration script: %s", err)
	}

	err = shared_kubernetes.RunPod(pgsqlFinalizeContainer, serverImage, flags.Image.PullPolicy, "/var/lib/uyuni-tools/"+scriptName, pgsqlFinalizeArgs...)
	if err != nil {
		return fmt.Errorf("Error running container %s: %s", pgsqlFinalizeContainer, err)
	}

	err = kubernetes.UyuniUpgrade(serverImage, flags.Image.PullPolicy, &flags.Helm, kubeconfig, fqdn, clusterInfos.Ingress)
	if err != nil {
		return fmt.Errorf("Cannot upgrade to image %s: %s", serverImage, err)
	}

	shared_kubernetes.WaitForDeployment(flags.Helm.Uyuni.Namespace, "uyuni", "uyuni")

	return nil
}
