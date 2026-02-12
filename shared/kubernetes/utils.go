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
var ServerFilter = fmt.Sprintf("-l%s=%s,%s=%s", AppLabel, ServerApp, ComponentLabel, "server")

// ServerFilter represents filter used to check proxy app.
const ProxyFilter = "-l" + AppLabel + "=" + ProxyApp

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
