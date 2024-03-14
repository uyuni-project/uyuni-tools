// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Connection contains information about how to connect to the server.
type Connection struct {
	backend          string
	command          string
	podName          string
	podmanContainer  string
	kubernetesFilter string
}

// Create a new connection object.
// The backend is either the command to use to connect to the container or the empty string.
//
// The empty strings means automatic detection of the backend where the uyuni container is running.
// podmanContainer is the name of a podman container to look for when detecting the command.
// kubernetesFilter is a filter parameter to use to match a pod.
func NewConnection(backend string, podmanContainer string, kubernetesFilter string) *Connection {
	cnx := Connection{backend: backend, podmanContainer: podmanContainer, kubernetesFilter: kubernetesFilter}

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
				err = fmt.Errorf("backend command not found in PATH: %s", c.backend)
			}
			c.command = c.backend
		case "":
			hasPodman := false
			hasKubectl := false

			// Check kubectl with a timeout in case the configured cluster is not responding
			_, err = exec.LookPath("kubectl")
			if err == nil {
				hasKubectl = true
				if out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "--request-timeout=30s", "get", "pod", c.kubernetesFilter, "-A", "-o=jsonpath={.items[*].metadata.name}"); err != nil {
					log.Info().Msg("kubectl not configured to connect to a cluster, ignoring")
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
					if checkErr := utils.RunCmd(bin, "inspect", c.podmanContainer, "--format", "{{.Name}}"); checkErr == nil {
						c.command = bin
						break
					}
				}
			}
			if c.command == "" {
				// Check for uyuni-server.service or helm release
				if hasPodman && podman.HasService("uyuni-server") {
					c.command = "podman"
				} else if hasKubectl {
					clusterInfos, err := kubernetes.CheckCluster()
					if err != nil {
						return c.command, err
					}
					if kubernetes.HasHelmRelease("uyuni", clusterInfos.GetKubeconfig()) {
						c.command = "kubectl"
					}
				}
			}
			if c.command == "" {
				err = fmt.Errorf("uyuni container is not accessible with one of podman, podman-remote or kubectl")
			}
		default:
			err = fmt.Errorf("unsupported backend %s", c.backend)
		}
	}
	return c.command, err
}

// GetPodName finds the name of the running pod.
func (c *Connection) GetPodName() (string, error) {
	var err error

	if c.podName == "" {
		command, cmdErr := c.GetCommand()
		if cmdErr != nil {
			log.Fatal().Err(cmdErr)
		}

		switch command {
		case "podman-remote":
			fallthrough
		case "podman":
			if out, _ := utils.RunCmdOutput(zerolog.DebugLevel, c.command, "ps", "-q", "-f", "name="+c.podmanContainer); len(out) == 0 {
				err = fmt.Errorf("container %s is not running on podman", c.podmanContainer)
			} else {
				c.podName = c.podmanContainer
			}
		case "kubectl":
			// We try the first item on purpose to make the command fail if not available
			podName, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", c.kubernetesFilter, "-A",
				"-o=jsonpath={.items[0].metadata.name}")
			if err == nil {
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
			return nil, fmt.Errorf("the container is not running, %s %s command not executed: %s",
				command, strings.Join(args, " "), err)
		}
	}

	cmd, cmdErr := c.GetCommand()
	if cmdErr != nil {
		return nil, cmdErr
	}

	cmdArgs := []string{"exec", c.podName}
	if cmd == "kubectl" {
		cmdArgs = append(cmdArgs, "-c", "uyuni", "--")
	}
	shellArgs := append([]string{command}, args...)
	cmdArgs = append(cmdArgs, shellArgs...)

	return utils.RunCmdOutput(zerolog.DebugLevel, cmd, cmdArgs...)
}

// WaitForServer waits at most 60s for multi-user systemd target to be reached.
func (c *Connection) WaitForServer() error {
	// Wait for the system to be up
	for i := 0; i < 60; i++ {
		podName, err := c.GetPodName()
		if err != nil {
			log.Fatal().Err(err)
		}

		args := []string{"exec", podName}
		command, err := c.GetCommand()
		if err != nil {
			log.Fatal().Err(err)
		}

		if command == "kubectl" {
			args = append(args, "--")
		}
		args = append(args, "systemctl", "is-active", "-q", "multi-user.target")
		output := utils.RunCmd(command, args...)
		isActive := output == nil

		if isActive {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return errors.New("server didn't start within 60s. Check for the service status")
}

// Copy transfers a file to or from the container.
// Prefix one of src or dst parameters with `server:` to designate the path is in the container
// user and group parameters are used to set the owner of a file transferred in the container.
func (c *Connection) Copy(src string, dst string, user string, group string) error {
	podName, err := c.GetPodName()
	if err != nil {
		return err
	}
	var commandArgs []string
	extraArgs := []string{}
	srcExpanded := strings.Replace(src, "server:", podName+":", 1)
	dstExpanded := strings.Replace(dst, "server:", podName+":", 1)

	command, err := c.GetCommand()
	if err != nil {
		return err
	}

	switch command {
	case "podman-remote":
		fallthrough
	case "podman":
		commandArgs = []string{"cp", srcExpanded, dstExpanded}
	case "kubectl":
		commandArgs = []string{"cp", "-c", "uyuni", srcExpanded, dstExpanded}
		extraArgs = []string{"-c", "uyuni", "--"}
	default:
		return fmt.Errorf("unknown container kind: %s", command)
	}

	if err := utils.RunCmdStdMapping(command, commandArgs...); err != nil {
		return err
	}

	if user != "" && strings.HasPrefix(dst, "server:") {
		execArgs := []string{"exec", podName}
		execArgs = append(execArgs, extraArgs...)
		owner := user
		if group != "" {
			owner = user + ":" + group
		}
		execArgs = append(execArgs, "chown", owner, strings.Replace(dst, "server:", "", 1))
		return utils.RunCmdStdMapping(command, execArgs...)
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
		log.Fatal().Msgf("Unknown container kind: %s\n", command)
	}

	if _, err := utils.RunCmdOutput(zerolog.DebugLevel, command, commandArgs...); err != nil {
		return false
	}
	return true
}

// ChoosePodmanOrKubernetes selects either the podman or the kubernetes function based on the backend.
// This function automatically detects the backend if compiled with kubernetes support and the backend flag is not passed.
func ChoosePodmanOrKubernetes[F interface{}](
	flags *pflag.FlagSet,
	podmanFn utils.CommandFunc[F],
	kubernetesFn utils.CommandFunc[F],
) (utils.CommandFunc[F], error) {
	backend := "podman"
	if utils.KubernetesBuilt {
		backend, _ = flags.GetString("backend")
	}

	cnx := NewConnection(backend, podman.ServerContainerName, kubernetes.ServerFilter)
	command, err := cnx.GetCommand()
	if err != nil {
		return nil, fmt.Errorf("failed to determine suitable backend")
	}
	switch command {
	case "podman":
		return podmanFn, nil
	case "kubectl":
		return kubernetesFn, nil
	}

	// Should never happen if the commands are the same than those handled in GetCommand()
	return nil, fmt.Errorf("no supported backend found")
}
