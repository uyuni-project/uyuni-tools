//go:build nok8s

package uninstall

const kubernetesBuilt = false

func uninstallForKubernetes(dryRun bool) {
}

func helmUninstall(kubeconfig string, deployment string, filter string, dryRun bool) string {
	return ""
}

const kubernetesHelp = ""
