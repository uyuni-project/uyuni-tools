package uninstall

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
)

func uninstallForKubernetes(globalFlags *types.GlobalFlags, dryRun bool) {
	clusterInfos := kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	// Uninstall uyuni
	namespace := helmUninstall(kubeconfig, "uyuni", "", dryRun)

	// Remove the remaining configmap and secrets
	if namespace != "" {
		if dryRun {
			log.Info().Msgf("Would run kubectl delete -n %s configmap uyuni-ca", namespace)
			log.Info().Msgf("Would run kubectl delete -n %s secret uyuni-ca uyuni-cert", namespace)
		} else {
			log.Info().Msgf("Running kubectl delete -n %s configmap uyuni-ca", namespace)
			if err := utils.RunRawCmd("kubectl", []string{"delete", "-n", namespace, "configmap", "uyuni-ca"}, false); err != nil {
				log.Info().Err(err).Msgf("Failed deleting config map")
			}

			log.Info().Msgf("Running kubectl delete -n %s secret uyuni-ca uyuni-cert", namespace)

			err := utils.RunRawCmd("kubectl", []string{"delete", "-n", namespace, "secret", "uyuni-ca", "uyuni-cert"}, false)
			if err != nil {
				log.Info().Err(err).Msgf("Failed deleting config map")
			}
		}
	}

	// Uninstall cert-manager if we installed it
	helmUninstall(kubeconfig, "cert-manager", "-linstalledby=uyuniadm", dryRun)

	// Remove the K3s Traefik config
	if clusterInfos.IsK3s() {
		kubernetes.UninstallK3sTraefikConfig(dryRun)
	}

	// Remove the rke2 nginx config
	if clusterInfos.IsRke2() {
		kubernetes.UninstallRke2NginxConfig(dryRun)
	}
}

func helmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) string {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].metadata.namespace}", deployment)
	args := []string{"get", "-A", "deploy", "-o", jsonpath}
	if filter != "" {
		args = append(args, filter)
	}

	out, err := exec.Command("kubectl", args...).Output()
	if err != nil {
		log.Info().Err(err).Msgf("Failed to find %s's namespace, skipping removal", deployment)
	}
	namespace := string(out)
	if namespace != "" {
		helmArgs := []string{}
		if kubeconfig != "" {
			helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
		}
		helmArgs = append(helmArgs, "uninstall", "-n", namespace, deployment)

		if dryRun {
			log.Info().Msgf("Would run helm %s", strings.Join(helmArgs, " "))
		} else {
			log.Info().Msgf("Uninstalling %s", deployment)
			message := "Failed to run helm " + strings.Join(helmArgs, " ")
			err := utils.RunRawCmd("helm", helmArgs, true)
			if err != nil {
				log.Fatal().Err(err).Msg(message)
			}
		}
	}
	return namespace
}
