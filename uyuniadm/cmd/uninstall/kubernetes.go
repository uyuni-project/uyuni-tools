package uninstall

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForKubernetes(globalFlags *types.GlobalFlags, dryRun bool) {
	// Uninstall uyuni
	namespace := helmUninstall("uyuni", "", dryRun, globalFlags.Verbose)

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
	helmUninstall("cert-manager", "-linstalledby=uyuniadm", dryRun, globalFlags.Verbose)
}

func helmUninstall(deployment string, filter string, dryRun bool, verbose bool) string {
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
		if dryRun {
			log.Printf("Would run helm uninstall %s\n", deployment)
		} else {
			log.Printf("Uninstalling %s\n", deployment)
			message := "Failed to run helm uninstall " + deployment
			utils.RunCmd("helm", []string{"uninstall", "-n", namespace, deployment}, message, verbose)
		}
	}
	return namespace
}
