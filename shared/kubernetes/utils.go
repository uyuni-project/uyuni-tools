// SPDX-FileCopyrightText: 2026 SUSE LLC
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
	"github.com/uyuni-project/uyuni-tools/shared/utils"
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

// ServerFilter represents filter used to check server app.
const ServerFilter = "-l" + AppLabel + "=" + ServerApp

// ServerFilter represents filter used to check proxy app.
const ProxyFilter = "-l" + AppLabel + "=" + ProxyApp

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
			ready, err := isDeploymentReady(namespace, deploymentName)
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

// isDeploymentReady returns true if a kubernetes deployment has at least one ready replica.
//
// An empty namespace means searching through all the namespaces.
func isDeploymentReady(namespace string, name string) (bool, error) {
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

	pods, err := getPodsForDeployment(namespace, name)
	if err != nil {
		return false, err
	}

	if err := hasAllPodsFailed(namespace, pods, name); err != nil {
		return false, err
	}

	return false, nil
}

func hasAllPodsFailed(namespace string, names []string, deployment string) error {
	failedPods := 0
	for _, podName := range names {
		if failed, err := isPodFailed(namespace, podName); err != nil {
			return err
		} else if failed {
			failedPods = failedPods + 1
		}
	}
	if len(names) > 0 && failedPods == len(names) {
		return fmt.Errorf(L("all the pods of %s deployment have a failure"), deployment)
	}
	return nil
}

func getPodsForDeployment(namespace string, name string) ([]string, error) {
	rs, err := getCurrentDeploymentReplicaSet(namespace, name)
	if err != nil {
		return []string{}, err
	}

	// Check if all replica set pods have failed to start
	return getPodsFromOwnerReference(namespace, rs)
}

func getCurrentDeploymentReplicaSet(namespace string, name string) (string, error) {
	// Get the replicasets matching the deployments and their revision as
	// Kubernetes doesn't remove the old replicasets after update.
	revisionPath := "{.metadata.annotations['deployment\\.kubernetes\\.io/revision']}"
	rsArgs := []string{
		"get", "rs", "-o",
		fmt.Sprintf(
			"jsonpath={range .items[?(@.metadata.ownerReferences[0].name=='%s')]}{.metadata.name},%s {end}",
			name, revisionPath,
		),
	}
	rsArgs = addNamespace(rsArgs, namespace)
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", rsArgs...)
	if err != nil {
		return "", utils.Errorf(err, L("failed to list ReplicaSets for deployment %s"), name)
	}
	replicasetsOut := strings.TrimSpace(string(out))
	// No replica, no deployment
	if replicasetsOut == "" {
		return "", nil
	}

	// Get the current deployment revision to look for
	out, err = runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "deploy", "-n", namespace, name,
		"-o", "jsonpath="+revisionPath,
	)
	if err != nil {
		return "", utils.Errorf(err, L("failed to get the %s deployment revision"), name)
	}
	revision := strings.TrimSpace(string(out))

	replicasets := strings.Split(replicasetsOut, " ")
	for _, rs := range replicasets {
		data := strings.SplitN(rs, ",", 2)
		if len(data) != 2 {
			return "", fmt.Errorf(L("invalid replicasset response: :%s"), replicasetsOut)
		}
		if data[1] == revision {
			return data[0], nil
		}
	}
	return "", nil
}

func getPodsFromOwnerReference(namespace string, owner string) ([]string, error) {
	jsonpath := fmt.Sprintf("jsonpath={.items[?(@.metadata.ownerReferences[0].name=='%s')].metadata.name}", owner)
	podArgs := []string{"get", "pod", "-o", jsonpath}
	podArgs = addNamespace(podArgs, namespace)
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", podArgs...)
	if err != nil {
		return []string{}, utils.Errorf(err, L("failed to find pods for owner reference %s"), owner)
	}

	outStr := strings.TrimSpace(string(out))

	pods := []string{}
	if outStr != "" {
		pods = strings.Split(outStr, " ")
	}
	return pods, nil
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

// ReplicasTo set the replicas for a deployment to the given value.
func ReplicasTo(namespace string, name string, replica uint) error {
	args := []string{"scale", "-n", namespace, "deploy", name, "--replicas", strconv.FormatUint(uint64(replica), 10)}
	log.Debug().Msgf("Setting replicas for deployment in %s to %d", name, replica)

	_, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", args...)
	if err != nil {
		return utils.Errorf(err, L("cannot run kubectl %s"), args)
	}

	if err := waitForReplicas(namespace, name, replica); err != nil {
		return err
	}

	log.Debug().Msgf("Replicas for %s deployment in %s are now %d", name, namespace, replica)
	return nil
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

func waitForReplicas(namespace string, name string, replicas uint) error {
	waitSeconds := 120
	log.Debug().Msgf("Checking replica for %s ready to %d", name, replicas)
	cmdArgs := []string{
		"get", "deploy", name, "-n", namespace, "-o", "jsonpath={.status.readyReplicas}", "--no-headers",
	}

	for i := 0; i < waitSeconds; i++ {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", cmdArgs...)
		if err != nil {
			return utils.Errorf(err, L("cannot execute %s"), strings.Join(cmdArgs, string(" ")))
		}
		outStr := strings.TrimSpace(string(out))
		var readyReplicas uint64
		if outStr != "" {
			var err error
			readyReplicas, err = strconv.ParseUint(outStr, 10, 8)
			if err != nil {
				return utils.Errorf(err, L("invalid replicas result"))
			}
		}
		if uint(readyReplicas) == replicas {
			return nil
		}
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
func GetPullPolicy(name string) string {
	policies := map[string]string{
		"always":       "Always",
		"never":        "Never",
		"ifnotpresent": "IfNotPresent",
	}
	policy := policies[strings.ToLower(name)]
	if policy == "" {
		log.Fatal().Msgf(L("%s is not a valid image pull policy value"), name)
	}
	return policy
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
