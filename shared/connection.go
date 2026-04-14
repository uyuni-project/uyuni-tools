// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
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
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var runner = utils.NewRunner

// SetRunner allows mocking the runner for tests.
func SetRunner(r func(command string, args ...string) types.Runner) {
	runner = r
}

// ResetRunner resets the runner to the default implementation.
func ResetRunner() {
	runner = utils.NewRunner
}

// Connection contains information about how to connect to the server.
type Connection struct {
	backend          string
	command          string
	podName          string
	kubernetesFilter string
	namespace        string
	container        string
	user             string
	systemd          podman.Systemd
}

// NewConnection creates a new connection object.
// The backend is either the command to use to connect to the container or the empty string.
//
// The empty strings means automatic detection of the backend where the uyuni container is running.
// container is the name of a container to look for when detecting the command.
// kubernetesFilter is a filter parameter to use to match a pod.
func NewConnection(backend string, container string, kubernetesFilter string) *Connection {
	return NewUserConnection(backend, container, kubernetesFilter, "")
}

// NewUserConnection creates a new connection object with a specific user to run commands as.
func NewUserConnection(backend string, container string, kubernetesFilter string, user string) *Connection {
	systemd := podman.NewSystemd()
	cnx := Connection{
		backend: backend, container: container, kubernetesFilter: kubernetesFilter, systemd: systemd, user: user,
	}

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
		case "host":
			c.command = "host"
		case "":
			hasPodman := false
			hasKubectl := false

			// Check kubectl with a timeout in case the configured cluster is not responding
			_, err = exec.LookPath("kubectl")
			if err == nil {
				hasKubectl = true
				if out, err := runner("kubectl", "--request-timeout=30s", "get", "deploy", c.kubernetesFilter,
					"-A", "-o=jsonpath={.items[*].metadata.name}",
				).Log(zerolog.DebugLevel).Spinner("").Exec(); err != nil {
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
					if _, checkErr := runner(bin, "inspect", c.container, "--format", "{{.Name}}").
						Spinner("").Exec(); checkErr == nil {
						c.command = bin
						break
					}
				}
			}
			if c.command == "" {
				// Check for uyuni-server.service or helm release
				if hasPodman && (c.systemd.HasService(podman.ServerService) || c.systemd.HasService(podman.ProxyService)) {
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
func (c *Connection) GetNamespace(appName string) (string, error) {
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

	// retrieving namespace from the first installed object we can find matching the filter.
	// This assumes that the server or proxy has been installed only in one namespace
	// with the current cluster credentials.
	out, err := runner("kubectl", "get", "all", "-A", c.kubernetesFilter, "-o",
		"jsonpath={.items[*].metadata.namespace}").Log(zerolog.DebugLevel).Spinner("").Exec()
	if err != nil {
		return "", utils.Errorf(err, L("failed to guest namespace"))
	}
	c.namespace = strings.TrimSpace(strings.Split(string(out), " ")[0])
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
			if out, _ := runner(c.command, "ps", "-q", "-f", "name="+c.container).
				Log(zerolog.DebugLevel).Spinner("").Exec(); len(out) == 0 {
				err = fmt.Errorf(L("container %s is not running on podman"), c.container)
			} else {
				log.Trace().Msgf("Found container ID '%s'", out)
				c.podName = c.container
			}
		case "kubectl":
			// We try the first item on purpose to make the command fail if not available
			if podName, _ := runner("kubectl", "get", "pod", c.kubernetesFilter, "-A",
				"-o=jsonpath={.items[0].metadata.name}").Log(zerolog.DebugLevel).Spinner("").Exec(); len(podName) == 0 {
				err = fmt.Errorf(L("container labeled %s is not running on kubectl"), c.kubernetesFilter)
			} else {
				c.podName = string(podName[:])
			}
		case "host":
			c.podName = "host"
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

	if cmd == "host" {
		if c.user != "" {
			fullCommand := quoteArgs(append([]string{command}, args...))
			return runner("su", "-", c.user, "-c", fullCommand).Log(zerolog.DebugLevel).Spinner("").Exec()
		}
		return runner(command, args...).Log(zerolog.DebugLevel).Spinner("").Exec()
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
	if c.user != "" {
		fullCommand := quoteArgs(shellArgs)
		shellArgs = []string{"su", "-", c.user, "-c", fullCommand}
	}
	cmdArgs = append(cmdArgs, shellArgs...)

	return runner(cmd, cmdArgs...).Log(zerolog.DebugLevel).Spinner("").Exec()
}

func quoteArgs(args []string) string {
	quotedArgs := make([]string, len(args))
	for i, arg := range args {
		quotedArgs[i] = "'" + strings.ReplaceAll(arg, "'", "'\\''") + "'"
	}
	return strings.Join(quotedArgs, " ")
}

// ExecScript runs the provided script inside the container.
func (c *Connection) ExecScript(script string) ([]byte, error) {
	tempFile, err := os.CreateTemp("", "uyuni-tools-script-*.sh")
	if err != nil {
		return nil, utils.Errorf(err, L("failed to create temporary file"))
	}
	defer os.Remove(tempFile.Name())

	if _, err = tempFile.WriteString(script); err != nil {
		return nil, utils.Errorf(err, L("failed to write script to temporary file"))
	}
	tempFile.Close()

	remotePath := fmt.Sprintf("/tmp/script-%d.sh", time.Now().UnixNano())

	// Copy localy created tempfile to container
	if err := c.Copy(tempFile.Name(), "server:"+remotePath, c.user, ""); err != nil {
		return nil, utils.Errorf(err, L("failed to copy script to container"))
	}

	defer func() {
		if _, err := c.Exec("rm", "-f", remotePath); err != nil {
			log.Debug().Err(err).Msgf("failed to remove %s", remotePath)
		}
	}()

	return c.Exec("bash", remotePath)
}

// Healthcheck runs healthcheck command inside the container.
func (c *Connection) Healthcheck() ([]byte, error) {
	if c.podName == "" {
		if _, err := c.GetPodName(); c.podName == "" {
			return nil, utils.Errorf(err, L("Healthcheck not executed"))
		}
	}

	cmd, cmdErr := c.GetCommand()
	if cmdErr != nil {
		return nil, cmdErr
	}

	if cmd == "host" {
		return nil, errors.New(L("healthcheck not supported on host"))
	}

	cmdArgs := []string{"healthcheck", "run", c.podName}

	return runner(cmd, cmdArgs...).Log(zerolog.DebugLevel).Spinner("").Exec()
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

		if command == "host" {
			return nil
		}

		if command == "kubectl" {
			args = append(args, "--")
		}
		args = append(args, "true")
		_, err = runner(command, args...).Spinner("").Exec()
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New(L("container didn't start within 10s."))
}

// WaitForHealthcheck waits at most 120s for healtcheck to succeed.
func (c *Connection) WaitForHealthcheck() error {
	// Wait for the system to be up
	for i := 0; i < 120; i++ {
		_, err := c.Healthcheck()
		if err != nil {
			log.Debug().Err(err)
			time.Sleep(1 * time.Second)
			continue
		}
		return nil
	}
	return errors.New(L("container didn't start within 120s. Check for the service status"))
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

	// Use connection user if not set explicitly
	if user == "" {
		user = c.user
	}

	switch command {
	case "podman-remote":
		fallthrough
	case "podman":
		commandArgs = []string{"cp", srcExpanded, dstExpanded}
	case "kubectl":
		commandArgs = []string{"cp", "-c", "uyuni", "-n", namespace, srcExpanded, dstExpanded}
		extraArgs = []string{"-c", "uyuni", "--"}
	case "host":
		srcExpanded = strings.Replace(src, "server:", "", 1)
		dstExpanded = strings.Replace(dst, "server:", "", 1)
		commandArgs = []string{srcExpanded, dstExpanded}
		command = "cp"
	default:
		return fmt.Errorf(L("unknown container kind: %s"), command)
	}

	if _, err := runner(command, commandArgs...).Log(zerolog.DebugLevel).StdMapping().Exec(); err != nil {
		return err
	}

	// File is already copied over, we need to drop server prefix
	dstPath := strings.Replace(dst, "server:", "", 1)
	if user != "" && (strings.HasPrefix(dst, "server:") || command == "host") {
		owner := user
		if group != "" {
			owner = user + ":" + group
		}

		execArgs := []string{"exec", podName}
		if command == "kubectl" {
			execArgs = append(execArgs, "-n", namespace)
		}
		if command == "host" {
			execArgs = []string{owner, dstPath}
			command = "chown"
		} else {
			execArgs = append(execArgs, extraArgs...)
			execArgs = append(execArgs, "chown", owner, dstPath)
		}

		_, err := runner(command, execArgs...).Log(zerolog.DebugLevel).StdMapping().Exec()
		return err
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
		namespace, err := c.GetNamespace("")
		if err != nil {
			log.Fatal().Err(err).Msg(L("failed to detect the namespace"))
		}
		commandArgs = append(commandArgs, "-n", namespace)
		commandArgs = append(commandArgs, "-c", "uyuni", "test", "-e", dstpath)
	case "host":
		return utils.FileExists(dstpath)
	default:
		log.Fatal().Msgf(L("unknown container kind: %s"), command)
	}

	if _, err := runner(command, commandArgs...).Log(zerolog.DebugLevel).Spinner("").Exec(); err != nil {
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
			if !utils.FileExists(pkiDir) {
				pkiDir = "/etc/ssl/certs" // OpenSSL fallback
			}
		}
	}
	hostPath := path.Join(pkiDir, fqdn+".crt")

	const containerCertPath = "server:/etc/pki/trust/anchors/LOCAL-RHN-ORG-TRUSTED-SSL-CERT"
	if err := c.Copy(containerCertPath, hostPath, "root", "root"); err != nil {
		return err
	}

	log.Info().Msg(L("Updating host trusted certificates"))
	if utils.CommandExists("update-ca-certificates") {
		_, err := runner("update-ca-certificates").Log(zerolog.DebugLevel).StdMapping().Exec()
		return err // openSUSE, Debian and Ubuntu
	} else if utils.CommandExists("update-ca-trust") {
		_, err := runner("update-ca-trust").Log(zerolog.DebugLevel).StdMapping().Exec()
		return err // RedHat
	} else if utils.CommandExists("trust") {
		_, err := runner("trust", "anchor", "--store", hostPath).Log(zerolog.DebugLevel).StdMapping().Exec()
		return err // Fallback
	}
	return errors.New(L("Unable to update host trusted certificates."))
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
	if runningBinary == "mgrpxy" {
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
func ChooseObjPodmanOrKubernetes[T any](systemd podman.Systemd, podmanOption T, kubernetesOption T) (T, error) {
	if systemd.HasService(podman.ServerService) || systemd.HasService(podman.ProxyService) {
		return podmanOption, nil
	} else if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		return kubernetesOption, nil
	}
	var res T
	return res, errors.New(L("failed to determine suitable backend"))
}

// RunSupportConfig will run supportconfig command on given connection.
func (c *Connection) RunSupportConfig(tmpDir string) ([]string, error) {
	var files []string
	const sourceBaseDir = "/var/log"

	containerName, err := c.GetPodName()
	if err != nil {
		return []string{}, err
	}

	// 10000 is what os.MkDirTemp uses
	const maxBatchNameAttempts = 10000
	batchName := ""
	sourceDir := ""
	for i := 0; i < maxBatchNameAttempts; i++ {
		suffix, err := utils.RandomHexString(4) // 8 hex chars
		if err != nil {
			return []string{}, fmt.Errorf(L("failed to generate supportconfig suffix: %w"), err)
		}

		candidateBatchName := "uyuni-server-container-" + suffix
		if !c.TestExistenceInPod(path.Join(sourceBaseDir, "scc_"+candidateBatchName)) {
			batchName = candidateBatchName
			sourceDir = path.Join(sourceBaseDir, "scc_"+candidateBatchName)
			break
		}
	}
	if batchName == "" {
		return []string{},
			fmt.Errorf(L("failed to generate unique supportconfig batch name after %d attempts"), maxBatchNameAttempts)
	}

	// Run supportconfig in the container if it's running
	log.Info().Msgf(L("Running supportconfig in %s"), containerName)
	if _, err = c.Exec("/sbin/supportconfig", "-B", batchName, "-t", sourceBaseDir); err != nil {
		/* do not return here.
		* supportconfig might return some error if some info is not generated
		* but we need to raise an error only if tarball is not generated.
		* In any case, show the error but as a warning and not as a failed run
		 */
		log.Warn().Err(err).Msg(L("Some parts of supportconfig were not successful"))
	}

	targetDir := path.Join(tmpDir, batchName)
	if err := c.Copy("server:"+sourceDir, targetDir, "", ""); err != nil {
		return []string{}, utils.Errorf(err, L("cannot copy support config"))
	}
	files = append(files, targetDir)

	// Remove the generated file in the container
	if _, err := c.Exec("rm", "-r", sourceDir); err != nil {
		return []string{}, utils.Errorf(err, L("failed to remove %s directory in the container"), sourceDir)
	}

	return files, nil
}
