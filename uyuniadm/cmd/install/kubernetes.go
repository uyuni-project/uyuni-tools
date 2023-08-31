package install

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/templates"
)

const HELM_APP_NAME = "uyuni"

func installForKubernetes(globalFlags *types.GlobalFlags, flags *InstallFlags, cmd *cobra.Command, args []string) {
	fqdn := args[0]

	// Check the kubernetes cluster setup
	clusterInfos := kubernetes.CheckCluster()

	// If installing on k3s, install the traefik helm config in manifests
	isK3s := clusterInfos.IsK3s()
	IsRke2 := clusterInfos.IsRke2()
	var kubeconfig string
	if isK3s {
		kubernetes.InstallK3sTraefikConfig()
		// If the user didn't provide a KUBECONFIG value or file, use the k3s default
		kubeconfigPath := os.ExpandEnv("${HOME}/.kube/config")
		if os.Getenv("KUBECONFIG") == "" || !utils.FileExists(kubeconfigPath) {
			kubeconfig = "/etc/rancher/k3s/k3s.yaml"
		}
	} else if IsRke2 {
		// Since even kubectl doesn't work without a trick on rke2, we assume the user has set kubeconfig
		kubernetes.InstallRke2NginxConfig(flags.Helm.Uyuni.Namespace)
	}

	if flags.Cert.UseExisting {
		// TODO Check that we have the expected secret and config in place
	} else {
		// Install cert-manager and a self-signed issuer ready for use
		installSslCertificates(globalFlags, flags, kubeconfig, fqdn)
	}

	// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	extractCaCertToConfig(globalFlags.Verbose)

	// Deploy the helm chart
	uyuniInstall(globalFlags, flags, kubeconfig, fqdn, clusterInfos.Ingress)

	// Wait for the pod to be started
	waitForDeployment(flags.Helm.Uyuni.Namespace, HELM_APP_NAME, "uyuni")
	utils.WaitForServer(globalFlags, "")

	// Create setup script + env variables and copy it to the container
	envs := map[string]string{
		"NO_SSL": "Y",
	}

	runSetup(globalFlags, flags, args[0], envs)
}

// Install cert-manager and its CRDs using helm in the cert-manager namespace if needed
// and then create a self-signed CA and issuers.
func installSslCertificates(globalFlags *types.GlobalFlags, flags *InstallFlags, kubeconfig, fqdn string) {
	// Install cert-manager if needed
	if !isDeploymentReady("", "cert-manager") {
		log.Println("Installing cert-manager")
		repo := ""
		chart := flags.Helm.CertManager.Chart
		version := flags.Helm.CertManager.Version
		namespace := flags.Helm.CertManager.Namespace

		args := []string{
			"--set", "installCRDs=true",
			"--set-json", "global.commonLabels={\"installedby\": \"uyuniadm\"}",
		}
		extraValues := flags.Helm.CertManager.Values
		if extraValues != "" {
			args = append(args, "-f", extraValues)
		}

		// Use upstream chart if nothing defined
		if chart == "" {
			repo = "https://charts.jetstack.io"
			chart = "cert-manager"
		}
		// The installedby label will be used to only uninstall what we installed
		helmInstall(globalFlags, kubeconfig, namespace, repo, "cert-manager", chart, version, args...)
	}

	// Wait for cert-manager to be ready
	waitForDeployment("", "cert-manager-webhook", "webhook")

	log.Println("Creating issuer for self signed SSL certificate authority")
	crdsDir, err := os.MkdirTemp("", "uyuniadm-*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %s\n", err)
	}
	defer os.RemoveAll(crdsDir)

	issuerPath := filepath.Join(crdsDir, "issuer.yaml")

	issuerData := templates.IssuerTemplateData{
		Country: flags.Cert.Country,
		State:   flags.Cert.State,
		City:    flags.Cert.City,
		Org:     flags.Cert.Org,
		OrgUnit: flags.Cert.OU,
		Email:   flags.Cert.Email,
		Fqdn:    fqdn,
	}

	if err = utils.WriteTemplateToFile(issuerData, issuerPath, 0500, true); err != nil {
		log.Fatalf("Failed to generate issuer definition: %s\n", err)
	}

	utils.RunCmd("kubectl", []string{"apply", "-f", filepath.Join(crdsDir, "issuer.yaml")},
		"Failed to create issuer", globalFlags.Verbose)

	// Wait for issuer to be ready
	for i := 0; i < 60; i++ {
		out, err := exec.Command("kubectl", "get", "-o=jsonpath={.status.conditions[*].type}",
			"issuer", "uyuni-ca-issuer").Output()
		if err == nil && string(out) == "Ready" {
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatalln("Issuer didn't turn ready after 60s")
}

func extractCaCertToConfig(verbose bool) {
	// TODO Replace with [trust-manager](https://cert-manager.io/docs/projects/trust-manager/) to automate this

	log.Println("Extracting CA certificate to a configmap")
	// Skip extracting if the configmap is already present
	out, err := exec.Command("kubectl", "get", "configmap", "uyuni-ca", "-o=jsonpath={.data.ca\\.crt}").Output()
	log.Printf("CA cert: %s\n", string(out))
	if err == nil && len(out) > 0 {
		log.Println("uyuni-ca configmap already existing, skipping extraction")
		return
	}

	out, err = exec.Command("kubectl", "get", "secret", "uyuni-ca", "-o=jsonpath={.data.ca\\.crt}").Output()
	if err != nil {
		log.Fatalf("Failed to get uyuni-ca certificate: %s\n", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(out))
	if err != nil {
		log.Fatalf("Failed to base64 decode CA certificate: %s", err)
	}

	message := fmt.Sprintf("Failed to create uyuni-ca config map from certificate: %s\n", err)
	valueArg := "--from-literal=ca.crt=" + string(decoded)
	utils.RunCmd("kubectl", []string{"create", "configmap", "uyuni-ca", valueArg}, message, verbose)
}

func uyuniInstall(globalFlags *types.GlobalFlags, flags *InstallFlags, kubeconfig string, fqdn string, ingress string) {
	log.Println("Installing Uyuni")

	// The issuer annotation is before the user's value to allow it to be overwritten for now.
	// Same for the guessed ingress: let the user override it in case we got it wrong.
	// TODO Parametrize the ca issuer value?
	helmParams := []string{
		"--set-json", "ingressSslAnnotations={\"cert-manager.io/issuer\": \"uyuni-ca-issuer\"}",
		"--set", "ingress=" + ingress,
	}

	extraValues := flags.Helm.Uyuni.Values
	if extraValues != "" {
		helmParams = append(helmParams, "-f", extraValues)
	}

	// The values computed from the command line need to be last to override what could be in the extras
	helmParams = append(helmParams,
		"--set", fmt.Sprintf("images.server=%s:%s", flags.Image.Name, flags.Image.Tag),
		"--set", "timezone="+flags.TZ,
		"--set", "fqdn="+fqdn)

	namespace := flags.Helm.Uyuni.Namespace
	chart := flags.Helm.Uyuni.Chart
	version := flags.Helm.Uyuni.Version
	helmInstall(globalFlags, kubeconfig, namespace, "", HELM_APP_NAME, chart, version, helmParams...)
}

// helmInstall runs helm install.
// If repo is not empty, the --repo parameter will be passed.
// If version is not empty, the --version parameter will be passed.
func helmInstall(globalFlags *types.GlobalFlags, kubeconfig string, namespace string, repo string, name string, chart string, version string, args ...string) {
	helmArgs := []string{
		"install",
		"-n", namespace,
		"--create-namespace",
		name,
		chart,
	}
	if kubeconfig != "" {
		helmArgs = append(helmArgs, "--kubeconfig", kubeconfig)
	}

	if repo != "" {
		helmArgs = append(helmArgs, "--repo", repo)
	}
	if version != "" {
		helmArgs = append(helmArgs, "--version", version)
	}
	helmArgs = append(helmArgs, args...)
	errorMessage := fmt.Sprintf("Failed to install helm chart %s in namespace %s", chart, namespace)

	utils.RunCmd("helm", helmArgs, errorMessage, globalFlags.Verbose)
}

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
		out, err := exec.Command("kubectl", cmdArgs...).Output()
		if err == nil {
			podName = string(out)
			break
		}
	}

	// We need to wait for the image to be pulled as this can add quite some time
	// Setting a timeout on this is very hard since it hightly depends on network speed and image size
	// List the Pulled events from the pod as we may not see the Pulling if the image was already downloaded
	waitForPulledImage(namespace, podName)

	log.Printf("Waiting for %s deployment to be ready in %s namespace\n", name, namespace)
	// Wait for a replica to be ready
	for i := 0; i < 60; i++ {
		// TODO Look for pod failures
		if isDeploymentReady(namespace, name) {
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatalf("Failed to find a ready replica for deployment %s in namespace %s after 60s\n", name, namespace)
}

func waitForPulledImage(namespace string, podName string) {
	log.Printf("Waiting for image of %s pod in %s namespace to be pulled\n", podName, namespace)
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
		out, err := exec.Command("kubectl", failedArgs...).Output()
		if err != nil {
			log.Fatalf("Failed to get failed events for pod %s: %s", podName, err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Failed to pull image") {
				log.Fatalln(err)
			}
		}

		// Has the image pull finished?
		out, err = exec.Command("kubectl", pulledArgs...).Output()
		if err != nil {
			log.Fatalf("Failed to get events for pod %s: %s\n", podName, err)
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

	out, err := exec.Command("kubectl", args...).Output()
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
