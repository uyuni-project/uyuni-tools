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
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const HELM_APP_NAME = "uyuni"

func installForKubernetes(viper *viper.Viper, globalFlags *types.GlobalFlags, cmd *cobra.Command, args []string) {
	fqdn := args[0]
	if viper.GetBool("cert.useexisting") {
		// TODO Check that we have the expected secret and config in place
	} else {
		// Install cert-manager and a self-signed issuer ready for use
		installSslCertificates(viper, fqdn, globalFlags)
	}

	// Extract the CA cert into uyuni-ca config map as the container shouldn't have the CA secret
	extractCaCertToConfig(globalFlags.Verbose)

	// Deploy the helm chart
	uyuniInstall(viper, fqdn, globalFlags)

	// Wait for the pod to be started
	waitForDeployment(viper.GetString("helm.namespace"), HELM_APP_NAME, "uyuni")
	utils.WaitForServer()

	// Create setup script + env variables and copy it to the container
	envs := map[string]string{
		"NO_SSL": "Y",
	}
	tmpFolder := generateSetupScript(viper, args[0], envs)
	defer os.RemoveAll(tmpFolder)

	utils.Copy(globalFlags, filepath.Join(tmpFolder, SETUP_NAME), "server:/tmp/setup.sh", "root", "root")

	// Run the setup script
	utils.Exec(globalFlags, false, false, []string{}, "/tmp/setup.sh")
}

// Install cert-manager and its CRDs using helm in the cert-manager namespace if needed
// and then create a self-signed CA and issuers.
func installSslCertificates(viper *viper.Viper, fqdn string, globalFlags *types.GlobalFlags) {
	// Install cert-manager if needed
	if !isDeploymentReady("", "cert-manager") {
		log.Println("Installing cert-manager")
		repo := ""
		chart := viper.GetString("helm.certmanager.chart")
		version := viper.GetString("helm.certmanager.version")
		namespace := viper.GetString("helm.certmanager.namespace")

		args := []string{
			"--set", "installCRDs=true",
			"--set-json", "global.commonLabels={\"installedby\": \"uyuniadm\"}",
		}
		extraValues := viper.GetString("helm.certmanager.values")
		if extraValues != "" {
			args = append(args, "-f", extraValues)
		}

		// Use upstream chart if nothing defined
		if chart == "" {
			repo = "https://charts.jetstack.io"
			chart = "cert-manager"
		}
		// The installedby label will be used to only uninstall what we installed
		helmInstall(globalFlags, namespace, repo, "cert-manager", chart, version, args...)
	}

	// Wait for cert-manager to be ready
	waitForDeployment("", "cert-manager-webhook", "webhook")

	// Deploy self-signed issuer
	const issuerTemplate = `apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-issuer
  namespace: default
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: uyuni-ca
  namespace: default
spec:
  isCA: true
  subject:
    countries: ["{{ .Country }}"]
    provinces: ["{{ .State }}"]
    localities: ["{{ .City }}"]
    organizations: ["{{ .Org }}"]
    organizationalUnits: ["{{ .OrgUnit }}"]
  emailAddresses:
    - {{ .Email }}
  commonName: {{ .Fqdn }}
  dnsNames:
    - {{ .Fqdn }}
  secretName: uyuni-ca
  privateKey:
    algorithm: ECDSA
    size: 256
  issuerRef:
    name: uyuni-issuer
    kind: Issuer
    group: cert-manager.io
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: uyuni-ca-issuer
  namespace: default
spec:
  ca:
    secretName:
      uyuni-ca
`

	log.Println("Creating issuer for self signed SSL certificate authority")
	crdsDir, err := os.MkdirTemp("", "uyuniadm-*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %s\n", err)
	}
	defer os.RemoveAll(crdsDir)

	issuerPath := filepath.Join(crdsDir, "issuer.yaml")
	file, err := os.OpenFile(issuerPath, os.O_WRONLY|os.O_CREATE, 0500)
	if err != nil {
		log.Fatalf("Fail to open %s file for writing: %s\n", issuerPath, err)
	}
	defer file.Close()

	model := struct {
		Country string
		State   string
		City    string
		Org     string
		OrgUnit string
		Email   string
		Fqdn    string
	}{
		Country: viper.GetString("cert.country"),
		State:   viper.GetString("cert.state"),
		City:    viper.GetString("cert.city"),
		Org:     viper.GetString("cert.org"),
		OrgUnit: viper.GetString("cert.ou"),
		Email:   viper.GetString("cert.email"),
		Fqdn:    fqdn,
	}

	t := template.Must(template.New("issuer").Parse(issuerTemplate))
	if err = t.Execute(file, model); err != nil {
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

func uyuniInstall(viper *viper.Viper, fqdn string, globalFlags *types.GlobalFlags) {
	log.Println("Installing Uyuni")

	// The issuer annotation is before the user's value to allow it to be overwritten for now.
	// TODO Parametrize the ca issuer value?
	helmParams := []string{"--set-json", "ingressSslAnnotations={\"cert-manager.io/issuer\": \"uyuni-ca-issuer\"}"}

	extraValues := viper.GetString("helm.uyuni.values")
	if extraValues != "" {
		helmParams = append(helmParams, "-f", extraValues)
	}

	// The values computed from the command line need to be last to override what could be in the extras
	helmParams = append(helmParams,
		"--set", fmt.Sprintf("images.server=%s:%s", viper.GetString("image"), viper.GetString("tag")),
		"--set", "timezone="+viper.GetString("tz"),
		"--set", "fqdn="+fqdn)

	namespace := viper.GetString("helm.uyuni.namespace")
	chart := viper.GetString("helm.uyuni.chart")
	version := viper.GetString("helm.uyuni.version")
	helmInstall(globalFlags, namespace, "", HELM_APP_NAME, chart, version, helmParams...)
}

// helmInstall runs helm install.
// If repo is not empty, the --repo parameter will be passed.
// If version is not empty, the --version parameter will be passed.
func helmInstall(globalFlags *types.GlobalFlags, namespace string, repo string, name string, chart string, version string, args ...string) {
	helmArgs := []string{
		"install",
		"-n", namespace,
		"--create-namespace",
		name,
		chart,
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
