package uninstall

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
)

func uninstallForKubernetes(globalFlags *types.GlobalFlags, dryRun bool) {
	clusterInfos := kubernetes.CheckCluster()
	kubeconfig := clusterInfos.GetKubeconfig()

	// Uninstall uyuni
	namespace := helmUninstall(kubeconfig, "uyuni", "", dryRun, globalFlags.Verbose)

	// Remove the remaining configmap and secrets
	if namespace != "" {
		if dryRun {
			log.Printf("Would run kubectl delete -n %s configmap uyuni-ca\n", namespace)
			log.Printf("Would run kubectl delete -n %s secret uyuni-ca uyuni-cert\n", namespace)
		} else {
			log.Printf("Running kubectl delete -n %s configmap uyuni-ca\n", namespace)
			if err := exec.Command("kubectl", "delete", "-n", namespace, "configmap", "uyuni-ca").Run(); err != nil {
				log.Printf("Failed deleting config map: %s\n", err)
			}

			log.Printf("Running kubectl delete -n %s secret uyuni-ca uyuni-cert\n", namespace)
			err := exec.Command("kubectl", "delete", "-n", namespace, "secret", "uyuni-ca", "uyuni-cert").Run()
			if err != nil {
				log.Printf("Failed deleting config map: %s\n", err)
			}
		}
	}

	// Uninstall cert-manager if we installed it
	helmUninstall(kubeconfig, "cert-manager", "-linstalledby=uyuniadm", dryRun, globalFlags.Verbose)

	// Remove the K3s Traefik config
	if clusterInfos.IsK3s() {
		kubernetes.UninstallK3sTraefikConfig(dryRun)
	}

	// Remove the rke2 nginx config
	if clusterInfos.IsRke2() {
		kubernetes.UninstallRke2NginxConfig(dryRun)
	}
}

func helmUninstall(kubeconfig string, deployment string, filter string, dryRun bool, verbose bool) string {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].metadata.namespace}", deployment)
	args := []string{"get", "-A", "deploy", "-o", jsonpath}
	if filter != "" {
		args = append(args, filter)
	}

	cmd := exec.Command("kubectl", args...)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to find %s's namespace, skipping removal: %s\n", deployment, err)
	}
	namespace := string(out)
	if namespace != "" {
		helmArgs := []string{}
		if kubeconfig != "" {
			helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
		}
		helmArgs = append(helmArgs, "uninstall", "-n", namespace, deployment)

		if dryRun {
			log.Printf("Would run helm %s\n", strings.Join(helmArgs, " "))
		} else {
			log.Printf("Uninstalling %s\n", deployment)
			message := "Failed to run helm " + strings.Join(helmArgs, " ")
			utils.RunCmd("helm", helmArgs, message, verbose)
		}
	}
	return namespace
}
