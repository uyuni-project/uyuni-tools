// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
)

const (
	// AppLabel is the app label name.
	AppLabel = "app.kubernetes.io/part-of"
	// ComponentLabel is the component label name.
	ComponentLabel = "app.kubernetes.io/component"
)

const (
	// ServerApp is the server app name.
	ServerApp = "uyuni"

	// ProxyApp is the proxy app name.
	ProxyApp = "uyuni-proxy"
)

const (
	// ServerComponent is the value of the component label for the server resources.
	ServerComponent = "server"
	// HubApiComponent is the value of the component label for the Hub API resources.
	HubAPIComponent = "hub-api"
	// CocoComponent is the value of the component label for the confidential computing attestation resources.
	CocoComponent = "coco"
)

// ServerFilter represents filter used to check server app.
const ServerFilter = "-l" + AppLabel + "=" + ServerApp

// ServerFilter represents filter used to check proxy app.
const ProxyFilter = "-l" + AppLabel + "=" + ProxyApp

// CaIssuerName is the name of the server CA issuer deployed if cert-manager is used.
const CaIssuerName = "uyuni-ca-issuer"

// GetLabels creates the label map with the app and component.
// The component label may be an empty string to skip it.
func GetLabels(app string, component string) map[string]string {
	labels := map[string]string{
		AppLabel: app,
	}
	if component != "" {
		labels[ComponentLabel] = component
	}
	return labels
}

// WaitForDeployment waits for a kubernetes deployment to have at least one replica.
func WaitForDeployments(namespace string, names ...string) error {
	log.Info().Msgf(
		NL("Waiting for %[1]s deployment to be ready in %[2]s namespace\n",
			"Waiting for %[1]s deployments to be ready in %[2]s namespace\n", len(names)),
		strings.Join(names, ", "), namespace)

	deploymentsStarting := names
	// Wait for ever for all deployments to be ready
	for len(deploymentsStarting) > 0 {
		starting := []string{}
		for _, deploymentName := range deploymentsStarting {
			ready, err := IsDeploymentReady(namespace, deploymentName)
			if err != nil {
				return err
			}
			if !ready {
				starting = append(starting, deploymentName)
			}
			deploymentsStarting = starting
		}
		if len(deploymentsStarting) > 0 {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// IsDeploymentReady returns true if a kubernetes deployment has at least one ready replica.
//
// An empty namespace means searching through all the namespaces.
func IsDeploymentReady(namespace string, name string) (bool, error) {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.name==\"%s\")].status.readyReplicas}", name)
	args := []string{"get", "-o", jsonpath, "deploy"}
	args = addNamespace(args, namespace)

	out, err := utils.RunCmdOutput(zerolog.TraceLevel, "kubectl", args...)
	// kubectl errors out if the deployment or namespace doesn't exist
	if err == nil {
		if replicas, _ := strconv.Atoi(string(out)); replicas > 0 {
			return true, nil
		}
	}

	// Search for the replica set matching the deployment
	rsArgs := []string{
		"get", "rs", "-o",
		fmt.Sprintf("jsonpath={.items[?(@.metadata.ownerReferences[0].name=='%s')].metadata.name}", name),
	}
	rsArgs = addNamespace(rsArgs, namespace)
	out, err = utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", rsArgs...)
	if err != nil {
		return false, utils.Errorf(err, L("failed to find ReplicaSet for deployment %s"), name)
	}
	rs := strings.TrimSpace(string(out))

	// Check if all replica set pods have failed to start
	jsonpath = fmt.Sprintf("jsonpath={.items[?(@.metadata.ownerReferences[0].name=='%s')].metadata.name}", rs)
	podArgs := []string{"get", "pod", "-o", jsonpath}
	podArgs = addNamespace(podArgs, namespace)
	out, err = utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", podArgs...)
	if err != nil {
		return false, utils.Errorf(err, L("failed to find pods for %s deployment"), name)
	}
	pods := strings.Split(string(out), " ")
	failedPods := 0
	for _, podName := range pods {
		if failed, err := isPodFailed(namespace, podName); err != nil {
			return false, err
		} else if failed {
			failedPods = failedPods + 1
		}
	}
	if failedPods == len(pods) {
		return false, fmt.Errorf(L("all the pods of %s deployment have a failure"), name)
	}

	return false, nil
}

// isPodFailed checks if any of the containers of the pod are in BackOff state.
//
// An empty namespace means searching through all the namespaces.
func isPodFailed(namespace string, name string) (bool, error) {
	// If a container failed to pull the image it status will have waiting.reason = ImagePullBackOff
	// If a container crashed its status will have waiting.reason = CrashLoopBackOff
	filter := fmt.Sprintf(".items[?(@.metadata.name==\"%s\")]", name)
	jsonpath := fmt.Sprintf("jsonpath={%[1]s.status.containerStatuses[*].state.waiting.reason}"+
		"{%[1]s.status.initContainerStatuses[*].state.waiting.reason}", filter)
	args := []string{"get", "pod", "-n", namespace, "-o", jsonpath}
	args = addNamespace(args, namespace)

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)
	if err != nil {
		return true, utils.Errorf(err, L("failed to get the status of %s pod"), name)
	}
	statuses := string(out)
	if strings.Contains(statuses, "CrashLoopBackOff") || strings.Contains(statuses, "ImagePullBackOff") {
		return true, nil
	}
	return false, nil
}

// DeploymentStatus represents the kubernetes deployment status.
type DeploymentStatus struct {
	AvailableReplicas int
	ReadyReplicas     int
	UpdatedReplicas   int
	Replicas          int
}

// GetDeploymentStatus returns the replicas status of the deployment.
func GetDeploymentStatus(namespace string, name string) (*DeploymentStatus, error) {
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "deploy", "-n", namespace,
		name, "-o", "jsonpath={.status}")
	if err != nil {
		return nil, err
	}

	var status DeploymentStatus
	if err = json.Unmarshal(out, &status); err != nil {
		return nil, utils.Errorf(err, L("failed to parse deployment status"))
	}
	return &status, nil
}

// ReplicasTo set the replica for an app to the given value.
// Scale the number of replicas of the server.
func ReplicasTo(namespace string, app string, replica uint) error {
	args := []string{"scale", "deploy", app, "--replicas"}
	log.Debug().Msgf("Setting replicas for pod in %s to %d", app, replica)
	args = append(args, fmt.Sprint(replica), "-n", namespace)

	_, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)
	if err != nil {
		return utils.Errorf(err, L("cannot run kubectl %s"), args)
	}

	pods, err := GetPods(namespace, "-l"+AppLabel+"="+app)
	if err != nil {
		return utils.Errorf(err, L("cannot get pods for %s"), app)
	}

	for _, pod := range pods {
		if len(pod) > 0 {
			err = waitForReplica(namespace, pod, replica)
			if err != nil {
				return utils.Errorf(err, L("replica to %d failed"), replica)
			}
		}
	}

	log.Debug().Msgf("Replicas for pod in %s are now %d", app, replica)

	return err
}

func isPodRunning(namespace string, podname string, filter string) (bool, error) {
	pods, err := GetPods(namespace, filter)
	if err != nil {
		return false, utils.Errorf(err, L("cannot check if pod %[1]s is running in app %[2]s"), podname, filter)
	}
	return utils.Contains(pods, podname), nil
}

// GetPods return the list of the pod given a filter.
func GetPods(namespace string, filter string) (pods []string, err error) {
	log.Debug().Msgf("Checking all pods for %s", filter)
	cmdArgs := []string{"get", "pods", "-n", namespace, filter, "--output=custom-columns=:.metadata.name", "--no-headers"}
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", cmdArgs...)
	if err != nil {
		return pods, utils.Errorf(err, L("cannot execute %s"), strings.Join(cmdArgs, string(" ")))
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, pod := range lines {
		pods = append(pods, strings.TrimSpace(pod))
	}
	log.Debug().Msgf("Pods in %s are %s", filter, pods)

	return pods, err
}

func waitForReplicaZero(namespace string, podname string) error {
	waitSeconds := 120
	cmdArgs := []string{"get", "pod", podname, "-n", namespace}

	for i := 0; i < waitSeconds; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", cmdArgs...)
		/* Assume that if the command return an error at the first iteration, it's because it failed,
		* next iteration because the pod was actually deleted
		 */
		if err != nil && i == 0 {
			return utils.Errorf(err, L("cannot get pod informations %s"), podname)
		}
		outStr := strings.TrimSuffix(string(out), "\n")
		if len(outStr) == 0 {
			log.Debug().Msgf("Pod %s has been deleted", podname)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf(L("cannot set replicas for %s to zero"), podname)
}

func waitForReplica(namespace string, podname string, replica uint) error {
	waitSeconds := 120
	log.Debug().Msgf("Checking replica for %s ready to %d", podname, replica)
	if replica == 0 {
		return waitForReplicaZero(namespace, podname)
	}
	cmdArgs := []string{
		"get", "pod", podname, "-n", namespace, "--output=custom-columns=STATUS:.status.phase", "--no-headers",
	}

	for i := 0; i < waitSeconds; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", cmdArgs...)
		if err != nil {
			return utils.Errorf(err, L("cannot execute %s"), strings.Join(cmdArgs, string(" ")))
		}
		outStr := strings.TrimSuffix(string(out), "\n")
		if string(outStr) == "Running" {
			log.Debug().Msgf("%s pod replica is now %d", podname, replica)
			break
		}
		log.Debug().Msgf("Pod %s replica is %s in %d seconds.", podname, string(out), i)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func addNamespace(args []string, namespace string) []string {
	if namespace != "" {
		args = append(args, "-n", namespace)
	} else {
		args = append(args, "-A")
	}
	return args
}

// GetPullPolicy returns the kubernetes PullPolicy value, if exists.
func GetPullPolicy(name string) core.PullPolicy {
	policies := map[string]core.PullPolicy{
		"always":       core.PullAlways,
		"never":        core.PullNever,
		"ifnotpresent": core.PullIfNotPresent,
	}
	policy := policies[strings.ToLower(name)]
	if policy == "" {
		log.Fatal().Msgf(L("%s is not a valid image pull policy value"), name)
	}
	return policy
}

// RunPod runs a pod, waiting for its execution and deleting it.
func RunPod(
	namespace string,
	podname string,
	filter string,
	image string,
	pullPolicy string,
	command string,
	override ...string,
) error {
	arguments := []string{
		"run", "--rm", "-n", namespace, "--attach", "--pod-running-timeout=3h", "--restart=Never", podname,
		"--image", image, "--image-pull-policy", pullPolicy, filter,
	}

	if len(override) > 0 {
		arguments = append(arguments, `--override-type=strategic`)
		for _, arg := range override {
			overrideParam := "--overrides=" + arg
			arguments = append(arguments, overrideParam)
		}
	}

	arguments = append(arguments, "--command", "--", command)
	err := utils.RunCmdStdMapping(zerolog.DebugLevel, "kubectl", arguments...)
	if err != nil {
		return utils.Errorf(err, PL("The first placeholder is a command",
			"cannot run %[1]s using image %[2]s"), command, image)
	}
	return nil
}

// DeletePod deletes a kubernetes pod named podname.
func DeletePod(namespace string, podname string, filter string) error {
	isRunning, err := isPodRunning(namespace, podname, filter)
	if err != nil {
		return utils.Errorf(err, L("cannot delete pod %s"), podname)
	}
	if !isRunning {
		log.Debug().Msgf("no need to delete pod %s because is not running", podname)
		return nil
	}
	arguments := []string{"delete", "pod", podname, "-n", namespace}
	_, err = utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", arguments...)
	if err != nil {
		return utils.Errorf(err, L("cannot delete pod %s"), podname)
	}
	return nil
}

// GetNode return the node where the app is running.
func GetNode(namespace string, filter string) (string, error) {
	nodeName := ""
	cmdArgs := []string{"get", "pod", "-n", namespace, filter, "-o", "jsonpath={.items[*].spec.nodeName}"}
	for i := 0; i < 60; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", cmdArgs...)
		if err == nil {
			nodeName = string(out)
			break
		}
	}
	if len(nodeName) > 0 {
		log.Debug().Msgf("Node name matching filter %s is: %s", filter, nodeName)
	} else {
		return "", fmt.Errorf(L("cannot find node name matching filter %s"), filter)
	}
	return nodeName, nil
}

// GenerateOverrideDeployment generate a JSON files represents the deployment information.
func GenerateOverrideDeployment(deployData types.Deployment) (string, error) {
	ret, err := json.Marshal(deployData)
	if err != nil {
		return "", utils.Errorf(err, L("cannot serialize pod definition override"))
	}
	return string(ret), nil
}

// GetRunningImage returns the image of containerName for the server running in the current system.
func GetRunningImage(containerName string) (string, error) {
	args := []string{
		"get", "pods", "-A", ServerFilter,
		"-o", "jsonpath={.items[0].spec.containers[?(@.name=='" + containerName + "')].image}",
	}
	image, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)

	log.Debug().Msgf("%[1]s container image is: %[2]s", containerName, image)
	if err != nil {
		return "", err
	}
	return strings.Trim(string(image), "\n"), nil
}
