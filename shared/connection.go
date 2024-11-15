// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Connection contains information about how to connect to the server.
type Connection struct {
	backend          string
	command          string
	podName          string
	kubernetesFilter string
	namespace        string
	container        string
}

// NewConnection creates a new connection object.
// The backend is either the command to use to connect to the container or the empty string.
//
// The empty strings means automatic detection of the backend where the uyuni container is running.
// container is the name of a container to look for when detecting the command.
// kubernetesFilter is a filter parameter to use to match a pod.
func NewConnection(backend string, container string, kubernetesFilter string) *Connection {
	cnx := Connection{backend: backend, container: container, kubernetesFilter: kubernetesFilter}

	return &cnx
}

// GetCommand validates or guesses the connection backend command.
func (c *Connection) GetCommand() (string, error) {
	var err error
	if c.command == "" {
		switch c.backend {
		case "podman":
			fallthrough
		case "podman-remote":
			fallthrough
		case "kubectl":
			if _, err = exec.LookPath(c.backend); err != nil {
				err = fmt.Errorf(L("backend command not found in PATH: %s"), c.backend)
			}
			c.command = c.backend
		case "":
			hasPodman := false
			hasKubectl := false

			// Check kubectl with a timeout in case the configured cluster is not responding
			_, err = exec.LookPath("kubectl")
			if err == nil {
				hasKubectl = true
				if out, err := utils.RunCmdOutput(
					zerolog.DebugLevel, "kubectl", "--request-timeout=30s", "get", "pod", c.kubernetesFilter, "-A",
					"-o=jsonpath={.items[*].metadata.name}",
				); err != nil {
					log.Info().Msg(L("kubectl not configured to connect to a cluster, ignoring"))
				} else if len(bytes.TrimSpace(out)) != 0 {
					c.command = "kubectl"
					return c.command, err
				}
			}

			// Search for other backends
			bins := []string{"podman", "podman-remote"}
			for _, bin := range bins {
				if _, err = exec.LookPath(bin); err == nil {
					hasPodman = true
					if checkErr := utils.RunCmd(bin, "inspect", c.container, "--format", "{{.Name}}"); checkErr == nil {
						c.command = bin
						break
					}
				}
			}
			if c.command == "" {
				// Check for uyuni-server.service or helm release
				if hasPodman && (podman.HasService(podman.ServerService) || podman.HasService(podman.ProxyService)) {
					c.command = "podman"
					return c.command, nil
				} else if hasKubectl {
					clusterInfos, err := kubernetes.CheckCluster()
					if err != nil {
						return c.command, err
					}
					kubeconfig := clusterInfos.GetKubeconfig()
					if kubernetes.HasHelmRelease("uyuni", kubeconfig) || kubernetes.HasHelmRelease("uyuni-proxy", kubeconfig) {
						c.command = "kubectl"
						return c.command, nil
					}
				}
			}
			if c.command == "" {
				err = errors.New(L("uyuni container is not accessible with one of podman, podman-remote or kubectl"))
			}
		default:
			err = fmt.Errorf(L("unsupported backend %s"), c.backend)
		}
	}
	return c.command, err
}

// GetNamespace finds the namespace of the running pod
// appName is the name of the application to look for, if not provided it will be guessed based on the filter.
// filters are additional filters to use to find the pod.
func (c *Connection) GetNamespace(appName string, filters ...string) (string, error) {
	// skip if namespace is already set
	if c.namespace != "" {
		return c.namespace, nil
	}

	command, cmdErr := c.GetCommand()
	if cmdErr != nil {
		return "", cmdErr
	}

	// skip if the command is not resolvable or does not target kubectl
	if command != "kubectl" {
		return c.namespace, nil
	}

	// if no appName is provided, we'll assume it based on its filter
	if appName == "" {
		switch c.kubernetesFilter {
		case kubernetes.ProxyFilter:
			appName = kubernetes.ProxyApp
		case kubernetes.ServerFilter:
			appName = kubernetes.ServerApp
		}

		if appName == "" {
			return "", errors.New(L("coundn't find app name"))
		}
	}

	// retrieving namespace from helm release
	clusterInfos, clusterInfosErr := kubernetes.CheckCluster()
	if clusterInfosErr != nil {
		return "", utils.Errorf(clusterInfosErr, L("failed to discover the cluster type"))
	}

	kubeconfig := clusterInfos.GetKubeconfig()
	if !kubernetes.HasHelmRelease(appName, kubeconfig) {
		return "", fmt.Errorf(L("no %s helm release installed on the cluster"), appName)
	}

	var namespaceErr error
	c.namespace, namespaceErr = extractNamespaceFromConfig(appName, kubeconfig, filters...)
	if namespaceErr != nil {
		return "", utils.Errorf(namespaceErr, L("failed to find the %s deployment namespace"), appName)
	}

	return c.namespace, nil
}

// GetPodName finds the name of the running pod.
func (c *Connection) GetPodName() (string, error) {
	var err error

	if c.podName == "" {
		command, cmdErr := c.GetCommand()
		if cmdErr != nil {
			return "", cmdErr
		}

		switch command {
		case "podman-remote":
			fallthrough
		case "podman":
			if out, _ := utils.RunCmdOutput(
				zerolog.DebugLevel, c.command, "ps", "-q", "-f", "name="+c.container,
			); len(out) == 0 {
				err = fmt.Errorf(L("container %s is not running on podman"), c.container)
			} else {
				log.Trace().Msgf("Found container ID '%s'", out)
				c.podName = c.container
			}
		case "kubectl":
			// We try the first item on purpose to make the command fail if not available
			if podName, _ := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", c.kubernetesFilter, "-A",
				"-o=jsonpath={.items[0].metadata.name}"); len(podName) == 0 {
				err = fmt.Errorf(L("container labeled %s is not running on kubectl"), c.kubernetesFilter)
			} else {
				c.podName = string(podName[:])
			}
		}
	}

	return c.podName, err
}

// Exec runs command inside the container within an sh shell.
func (c *Connection) Exec(command string, args ...string) ([]byte, error) {
	if c.podName == "" {
		if _, err := c.GetPodName(); c.podName == "" {
			commandStr := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
			return nil, utils.Errorf(err, L("%s command not executed:"), commandStr)
		}
	}

	cmd, cmdErr := c.GetCommand()
	if cmdErr != nil {
		return nil, cmdErr
	}

	cmdArgs := []string{"exec", c.podName}
	if cmd == "kubectl" {
		if _, err := c.GetNamespace(""); c.namespace == "" {
			return nil, utils.Errorf(err, L("failed to retrieve namespace "))
		}

		if c.container == "" {
			c.container = "uyuni"
		}

		cmdArgs = append(cmdArgs, "-n", c.namespace, "-c", c.container, "--")
	}
	shellArgs := append([]string{command}, args...)
	cmdArgs = append(cmdArgs, shellArgs...)

	return utils.RunCmdOutput(zerolog.DebugLevel, cmd, cmdArgs...)
}

// WaitForContainer waits up to 10 sec for the container to appear.
func (c *Connection) WaitForContainer() error {
	for i := 0; i < 10; i++ {
		podName, err := c.GetPodName()
		if err != nil {
			log.Debug().Err(err)
			time.Sleep(1 * time.Second)
			continue
		}
		args := []string{"exec", podName}
		command, err := c.GetCommand()
		if err != nil {
			return err
		}

		if command == "kubectl" {
			args = append(args, "--")
		}
		args = append(args, "true")
		err = utils.RunCmd(command, args...)
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New(L("container didn't start within 10s."))
}

// WaitForServer waits at most 60s for multi-user systemd target to be reached.
func (c *Connection) WaitForServer() error {
	// Wait for the system to be up
	for i := 0; i < 60; i++ {
		podName, err := c.GetPodName()
		if err != nil {
			log.Debug().Err(err)
			time.Sleep(1 * time.Second)
			continue
		}

		namespace, err := c.GetNamespace("")
		if err != nil {
			return err
		}

		args := []string{"exec", podName}
		command, err := c.GetCommand()
		if err != nil {
			return err
		}

		if command == "kubectl" {
			args = append(args, "-n", namespace, "--")
		}
		args = append(args, "systemctl", "is-active", "-q", "multi-user.target")
		output := utils.RunCmd(command, args...)
		isActive := output == nil

		if isActive {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New(L("server didn't start within 60s. Check for the service status"))
}

// Copy transfers a file to or from the container.
// Prefix one of src or dst parameters with `server:` to designate the path is in the container
// user and group parameters are used to set the owner of a file transferred in the container.
func (c *Connection) Copy(src string, dst string, user string, group string) error {
	podName, err := c.GetPodName()
	if err != nil {
		return err
	}

	command, err := c.GetCommand()
	if err != nil {
		return err
	}

	var namespace, namespacePrefix = "", ""
	if command == "kubectl" {
		namespace, err = c.GetNamespace("")
		if err != nil {
			return err
		}
		namespacePrefix = namespace + "/"
	}

	var commandArgs []string
	extraArgs := []string{}
	srcExpanded := strings.Replace(src, "server:", namespacePrefix+podName+":", 1)
	dstExpanded := strings.Replace(dst, "server:", namespacePrefix+podName+":", 1)

	switch command {
	case "podman-remote":
		fallthrough
	case "podman":
		commandArgs = []string{"cp", srcExpanded, dstExpanded}
	case "kubectl":
		commandArgs = []string{"cp", "-c", "uyuni", "-n", namespace, srcExpanded, dstExpanded}
		extraArgs = []string{"-c", "uyuni", "--"}
	default:
		return fmt.Errorf(L("unknown container kind: %s"), command)
	}

	if err := utils.RunCmdStdMapping(zerolog.DebugLevel, command, commandArgs...); err != nil {
		return err
	}

	if user != "" && strings.HasPrefix(dst, "server:") {
		execArgs := []string{"exec", podName}
		if command == "kubectl" {
			execArgs = append(execArgs, "-n", namespace)
		}
		execArgs = append(execArgs, extraArgs...)
		owner := user
		if group != "" {
			owner = user + ":" + group
		}
		execArgs = append(execArgs, "chown", owner, strings.Replace(dst, "server:", "", 1))
		return utils.RunCmdStdMapping(zerolog.DebugLevel, command, execArgs...)
	}
	return nil
}

// TestExistenceInPod returns true if dstpath exists in the pod.
func (c *Connection) TestExistenceInPod(dstpath string) bool {
	podName, err := c.GetPodName()
	if err != nil {
		log.Fatal().Err(err)
	}
	commandArgs := []string{"exec", podName}

	command, err := c.GetCommand()
	if err != nil {
		log.Fatal().Err(err)
	}

	switch command {
	case "podman":
		commandArgs = append(commandArgs, "test", "-e", dstpath)
	case "kubectl":
		commandArgs = append(commandArgs, "-c", "uyuni", "test", "-e", dstpath)
	default:
		log.Fatal().Msgf(L("unknown container kind: %s"), command)
	}

	if _, err := utils.RunCmdOutput(zerolog.DebugLevel, command, commandArgs...); err != nil {
		return false
	}
	return true
}

// CopyCaCertificate copies the server SSL CA certificate to the host with fqdn as the name of the created file.
func (c *Connection) CopyCaCertificate(fqdn string) error {
	log.Info().Msg(L("Copying the SSL CA certificate to the host"))

	pkiDir := "/etc/pki/trust/anchors/"
	if !utils.FileExists(pkiDir) {
		pkiDir = "/etc/pki/ca-trust/source/anchors" // RedHat
		if !utils.FileExists(pkiDir) {
			pkiDir = "/usr/local/share/ca-certificates" // Debian and Ubuntu
		}
	}
	hostPath := path.Join(pkiDir, fqdn+".crt")

	const containerCertPath = "server:/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT"
	if err := c.Copy(containerCertPath, hostPath, "root", "root"); err != nil {
		return err
	}

	log.Info().Msg(L("Updating host trusted certificates"))
	return utils.RunCmdStdMapping(zerolog.DebugLevel, "update-ca-certificates")
}

// ChoosePodmanOrKubernetes selects either the podman or the kubernetes function based on the backend.
//
// This function automatically detects the backend if compiled with kubernetes support
// and the backend flag is not passed.
func ChoosePodmanOrKubernetes[F interface{}](
	flags *pflag.FlagSet,
	podmanFn utils.CommandFunc[F],
	kubernetesFn utils.CommandFunc[F],
) (utils.CommandFunc[F], error) {
	backend := "podman"
	runningBinary := filepath.Base(os.Args[0])
	if utils.KubernetesBuilt || runningBinary == "mgrpxy" {
		backend, _ = flags.GetString("backend")
	}

	cnx := NewConnection(backend, podman.ServerContainerName, kubernetes.ServerFilter)
	return chooseBackend(cnx, podmanFn, kubernetesFn)
}

// ChooseProxyPodmanOrKubernetes selects either the podman or the kubernetes function based on the proxy backend.
func ChooseProxyPodmanOrKubernetes[F interface{}](
	flags *pflag.FlagSet,
	podmanFn utils.CommandFunc[F],
	kubernetesFn utils.CommandFunc[F],
) (utils.CommandFunc[F], error) {
	backend, _ := flags.GetString("backend")

	cnx := NewConnection(backend, podman.ProxyContainerNames[0], kubernetes.ProxyFilter)
	return chooseBackend(cnx, podmanFn, kubernetesFn)
}

func chooseBackend[F interface{}](
	cnx *Connection,
	podmanFn utils.CommandFunc[F],
	kubernetesFn utils.CommandFunc[F],
) (utils.CommandFunc[F], error) {
	command, err := cnx.GetCommand()
	if err != nil {
		return nil, errors.New(L("failed to determine suitable backend"))
	}
	switch command {
	case "podman":
		return podmanFn, nil
	case "kubectl":
		return kubernetesFn, nil
	}

	// Should never happen if the commands are the same than those handled in GetCommand()
	return nil, errors.New(L("no supported backend found"))
}

// ChooseObjPodmanOrKubernetes returns an artibraty object depending if podman or the kubernetes is installed.
func ChooseObjPodmanOrKubernetes[T any](podmanOption T, kubernetesOption T) (T, error) {
	if podman.HasService(podman.ServerService) || podman.HasService(podman.ProxyService) {
		return podmanOption, nil
	} else if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		return kubernetesOption, nil
	}
	var res T
	return res, errors.New(L("failed to determine suitable backend"))
}

// RunSupportConfig will run supportconfig command on given connection.
func (c *Connection) RunSupportConfig(tmpDir string) ([]string, error) {
	var containerTarball string
	var files []string
	extensions := []string{"", ".md5"}
	containerName, err := c.GetPodName()
	if err != nil {
		return []string{}, err
	}

	// Run supportconfig in the container if it's running
	log.Info().Msgf(L("Running supportconfig in  %s"), containerName)
	out, err := c.Exec("supportconfig")
	if err != nil {
		return []string{}, errors.New(L("failed to run supportconfig"))
	}
	tarballPath := utils.GetSupportConfigPath(string(out))
	if tarballPath == "" {
		return []string{}, errors.New(L("failed to find container supportconfig tarball from command output"))
	}

	for _, ext := range extensions {
		containerTarball = path.Join(tmpDir, containerName+"-supportconfig.txz"+ext)
		if err := c.Copy("server:"+tarballPath+ext, containerTarball, "", ""); err != nil {
			return []string{}, utils.Errorf(err, L("cannot copy tarball"))
		}
		files = append(files, containerTarball)

		// Remove the generated file in the container
		if _, err := c.Exec("rm", tarballPath+ext); err != nil {
			return []string{}, utils.Errorf(err, L("failed to remove %s file in the container"), tarballPath+ext)
		}
	}
	return files, nil
}

// extractNamespaceFromConfig extracts the namespace of a given application
// from the Helm release information.
func extractNamespaceFromConfig(appName string, kubeconfig string, filters ...string) (string, error) {
	args := []string{}
	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}
	args = append(args, "list", "-aA", "-f", appName, "-o", "json")
	args = append(args, filters...)

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "helm", args...)
	if err != nil {
		return "", utils.Errorf(err, L("failed to detect %s's namespace using helm"), appName)
	}

	var data []releaseInfo
	if err = json.Unmarshal(out, &data); err != nil {
		return "", utils.Errorf(err, L("helm provided an invalid JSON output"))
	}

	if len(data) == 1 {
		return data[0].Namespace, nil
	}
	return "", errors.New(L("found no or more than one deployment"))
}

type releaseInfo struct {
	Namespace string `mapstructure:"namespace"`
}
