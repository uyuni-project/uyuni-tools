// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ClusterInfos represent cluster information.
type ClusterInfos struct {
	KubeletVersion string
	Ingress        string
}

// IsK3s is true if it's a K3s Cluster.
func (infos ClusterInfos) IsK3s() bool {
	return strings.Contains(infos.KubeletVersion, "k3s")
}

// IsRKE2 is true if it's a RKE2 Cluster.
func (infos ClusterInfos) IsRke2() bool {
	return strings.Contains(infos.KubeletVersion, "rke2")
}

// GetKubeconfig returns the path to the default kubeconfig file or "" if none.
func (infos ClusterInfos) GetKubeconfig() string {
	var kubeconfig string
	if infos.IsK3s() {
		// If the user didn't provide a KUBECONFIG value or file, use the k3s default
		kubeconfigPath := os.ExpandEnv("${HOME}/.kube/config")
		if os.Getenv("KUBECONFIG") == "" || !utils.FileExists(kubeconfigPath) {
			kubeconfig = "/etc/rancher/k3s/k3s.yaml"
		}
	}
	// Since even kubectl doesn't work without a trick on rke2, we assume the user has set kubeconfig
	return kubeconfig
}

// CheckCluster return cluster information.
func CheckCluster() (*ClusterInfos, error) {
	// Get the kubelet version
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "node",
		"-o", "jsonpath={.items[0].status.nodeInfo.kubeletVersion}")
	if err != nil {
		return nil, utils.Errorf(err, L("failed to get kubelet version"))
	}

	var infos ClusterInfos
	infos.KubeletVersion = string(out)
	infos.Ingress, err = guessIngress()
	if err != nil {
		return nil, err
	}

	return &infos, nil
}

func guessIngress() (string, error) {
	// Check for a traefik resource
	err := utils.RunCmd("kubectl", "explain", "ingressroutetcp")
	if err == nil {
		return "traefik", nil
	} else {
		log.Debug().Err(err).Msg("No ingressroutetcp resource deployed")
	}

	// Look for a pod running the nginx-ingress-controller: there is no other common way to find out
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", "-A",
		"-o", "jsonpath={range .items[*]}{.spec.containers[*].args[0]}{.spec.containers[*].command}{end}")
	if err != nil {
		return "", utils.Errorf(err, L("failed to get pod commands to look for nginx controller"))
	}

	const nginxController = "/nginx-ingress-controller"
	if strings.Contains(string(out), nginxController) {
		return "nginx", nil
	}

	return "", nil
}

// Restart restarts the pod.
func Restart(namespace string, app string) error {
	if err := Stop(namespace, app); err != nil {
		return utils.Errorf(err, L("cannot stop %s"), app)
	}
	return Start(namespace, app)
}

// Start starts the pod.
func Start(namespace string, app string) error {
	// if something is running, we don't need to set replicas to 1
	if _, err := GetNode(namespace, "-lapp="+app); err != nil {
		return ReplicasTo(namespace, app, 1)
	}
	log.Debug().Msgf("Already running")
	return nil
}

// Stop stop the pod.
func Stop(namespace string, app string) error {
	return ReplicasTo(namespace, app, 0)
}

func get(component string, componentName string, args ...string) ([]byte, error) {
	kubectlArgs := []string{
		"get",
		component,
		componentName,
	}

	kubectlArgs = append(kubectlArgs, args...)

	output, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", kubectlArgs...)
	if err != nil {
		return []byte{}, err
	}
	return output, nil
}

// GetConfigMap returns the value of a given config map.
func GetConfigMap(configMapName string, filter string) (string, error) {
	out, err := get("configMap", configMapName, filter)
	if err != nil {
		return "", utils.Errorf(err, L("failed to run kubectl get configMap %[1]s %[2]s"), configMapName, filter)
	}

	return string(out), nil
}

// GetSecret returns the value of a given secret.
func GetSecret(secretName string, filter string) (string, error) {
	out, err := get("secret", secretName, filter)
	if err != nil {
		return "", utils.Errorf(err, L("failed to run kubectl get secret %[1]s %[2]s"), secretName, filter)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		return "", utils.Errorf(err, L("Failed to base64 decode secret %s"), secretName)
	}

	return string(decoded), nil
}

// CreateDockerSecret creates a secret of docker type to authenticate registries.
func CreateDockerSecret(namespace string, name string, registry string, username string, password string) error {
	authString := fmt.Sprintf("%s:%s", username, password)
	auth := base64.StdEncoding.EncodeToString([]byte(authString))
	configjson := fmt.Sprintf(
		`{"auths": {"%s": {"username": "%s", "password": "%s", "auth": "%s"}}}`,
		registry, username, password, auth,
	)

	secret := fmt.Sprintf(`
apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  namespace: %s
  name: %s
data:
  .dockerconfigjson: %s
`, namespace, name, base64.StdEncoding.EncodeToString([]byte(configjson)))

	tempDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Run the job
	definitionPath := path.Join(tempDir, "definition.yaml")
	if err := os.WriteFile(definitionPath, []byte(secret), 0600); err != nil {
		return utils.Errorf(err, L("failed to write %s secret definition file"), name)
	}

	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, "kubectl", "apply", "-f", definitionPath); err != nil {
		return utils.Errorf(err, L("failed to define %s secret"), name)
	}
	return nil
}

// GetDeploymentImagePullSecret returns the name of the image pull secret of a deployment.
//
// This assumes only one secret is defined on the deployment.
func GetDeploymentImagePullSecret(namespace string, filter string) (string, error) {
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "deploy", "-n", namespace, filter,
		"-o", "jsonpath={.items[*].spec.template.spec.imagePullSecrets[*].name}",
	)
	if err != nil {
		return "", utils.Errorf(err, L("failed to get deployment image pull secret"))
	}

	return strings.TrimSpace(string(out)), nil
}
