// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

// IsRke2 is true if it's a RKE2 Cluster.
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
	}
	log.Debug().Err(err).Msg("No ingressroutetcp resource deployed")

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
	if _, err := GetNode(namespace, "-l"+AppLabel+"="+app); err != nil {
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

// createDockerSecret creates a secret of docker type to authenticate registries.
func createDockerSecret(
	namespace string,
	name string,
	registry string,
	username string,
	password string,
	appLabel string,
) error {
	authString := fmt.Sprintf("%s:%s", username, password)
	auth := base64.StdEncoding.EncodeToString([]byte(authString))
	configjson := fmt.Sprintf(
		`{"auths": {"%s": {"username": "%s", "password": "%s", "auth": "%s"}}}`,
		registry, username, password, auth,
	)

	secret := core.Secret{
		TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels:    GetLabels(appLabel, ""),
		},
		// It seems serializing this object automatically transforms the secrets to base64.
		Data: map[string][]byte{
			".dockerconfigjson": []byte(configjson),
		},
		Type: core.SecretTypeDockerConfigJson,
	}
	return Apply([]runtime.Object{&secret}, fmt.Sprintf(L("failed to create the %s docker secret"), name))
}

// AddRegistry creates a secret holding the registry credentials and adds it to the helm args.
func AddRegistry(helmArgs []string, namespace string, registry *types.Registry, appLabel string) ([]string, error) {
	secret, err := GetRegistrySecret(namespace, registry, appLabel)
	if secret != "" {
		helmArgs = append(helmArgs, secret)
	}
	return helmArgs, err
}

// GetRegistrySecret creates a docker secret holding the registry credentials and returns the secret name.
func GetRegistrySecret(namespace string, registry *types.Registry, appLabel string) (string, error) {
	const secretName = "registry-credentials"

	// Return the existing secret if any.
	out, err := runCmdOutput(zerolog.DebugLevel, "kubectl", "get", "-n", namespace, "secret", secretName, "-o", "name")
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return secretName, nil
	}

	// Create the secret if registry user and password are passed.
	if registry.User != "" && registry.Password != "" {
		if err := createDockerSecret(
			namespace, secretName, registry.Host, registry.User, registry.Password, appLabel,
		); err != nil {
			return "", err
		}
		return secretName, nil
	}
	return "", nil
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

// HasResource checks if a resource is available on the cluster.
func HasResource(name string) bool {
	if err := utils.RunCmd("kubectl", "explain", name); err != nil {
		return false
	}
	return true
}
