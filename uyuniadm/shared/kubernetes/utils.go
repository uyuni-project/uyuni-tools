package kubernetes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// waitForDeployment waits at most 60s for a kubernetes deployment to have at least one replica.
// See [isDeploymentReady] for more details.
func waitForDeployment(namespace string, name string, appName string) {
	// Find the name of a replica pod
	// Using the app label is a shortcut, not the 100% acurate way to get from deployment to pod
	podName := ""
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.labels.app==\"%s\")].metadata.name}", appName)
	cmdArgs := []string{"get", "pod", "-o", jsonpath}
	cmdArgs = addNamespace(cmdArgs, namespace)

	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput("kubectl", cmdArgs...)
		if err == nil {
			podName = string(out)
			break
		}
	}

	// We need to wait for the image to be pulled as this can add quite some time
	// Setting a timeout on this is very hard since it hightly depends on network speed and image size
	// List the Pulled events from the pod as we may not see the Pulling if the image was already downloaded
	waitForPulledImage(namespace, podName)

	log.Info().Msgf("Waiting for %s deployment to be ready in %s namespace\n", name, namespace)
	// Wait for a replica to be ready
	for i := 0; i < 60; i++ {
		// TODO Look for pod failures
		if isDeploymentReady(namespace, name) {
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatal().Msgf("Failed to find a ready replica for deployment %s in namespace %s after 60s", name, namespace)
}

func waitForPulledImage(namespace string, podName string) {
	log.Info().Msgf("Waiting for image of %s pod in %s namespace to be pulled", podName, namespace)
	pulledArgs := []string{"get", "event",
		"-o", "jsonpath={.items[?(@.reason==\"Pulled\")].message}",
		"--field-selector", "involvedObject.name=" + podName}
	pulledArgs = addNamespace(pulledArgs, namespace)

	failedArgs := []string{"get", "event",
		"-o", "jsonpath={range .items[?(@.reason==\"Failed\")]}{.message}{\"\\n\"}{end}",
		"--field-selector", "involvedObject.name=" + podName}
	failedArgs = addNamespace(failedArgs, namespace)
	for {
		// Look for events indicating an image pull issue
		out, err := utils.RunCmdOutput("kubectl", failedArgs...)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to get failed events for pod %s", podName)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Failed to pull image") {
				log.Fatal().Err(err).Msg("Failed to pull image")
			}
		}

		// Has the image pull finished?
		out, err = utils.RunCmdOutput("kubectl", pulledArgs...)
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to get events for pod %s", podName)
		}
		if len(out) > 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// isDeploymentReady returns true if a kubernetes deployment has at least one ready replica.
// The name can also be a filter parameter like -lapp=uyuni.
// An empty namespace means searching through all the namespaces.
func isDeploymentReady(namespace string, name string) bool {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].status.readyReplicas}", name)
	args := []string{"get", "-o", jsonpath, "deploy"}
	args = addNamespace(args, namespace)

	out, err := utils.RunCmdOutput("kubectl", args...)
	// kubectl errors out if the deployment or namespace doesn't exist
	if err == nil {
		if replicas, _ := strconv.Atoi(string(out)); replicas > 0 {
			return true
		}
	}
	return false
}

func addNamespace(args []string, namespace string) []string {
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	return args
}

func uninstallFile(path string, dryRun bool) {
	if utils.FileExists(path) {
		if dryRun {
			log.Info().Msgf("Would remove file %s", path)
		} else {
			log.Info().Msgf("Removing file %s", path)
			if err := os.Remove(path); err != nil {
				log.Info().Err(err).Msgf("Failed to remove file %s", path)
			}
		}
	}
}
