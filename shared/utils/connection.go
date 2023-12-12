// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
			// Check kubectl with a timeout in case the configured cluster is not responding
			_, err = exec.LookPath("kubectl")
			if err == nil {
				if out, err := RunCmdOutput(zerolog.DebugLevel, "kubectl", "--request-timeout=30s", "get", "pod", c.kubernetesFilter, "-A", "-o=jsonpath={.items[*].metadata.name}"); err != nil {
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
					if checkErr := RunCmd(bin, "inspect", c.podmanContainer, "--format", "{{.Name}}"); checkErr == nil {
						c.command = bin
						break
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
			if out, _ := RunCmdOutput(zerolog.DebugLevel, c.command, "ps", "-q", "-f", "name="+c.podmanContainer); len(out) == 0 {
				err = fmt.Errorf("container %s is not running on podman", c.podmanContainer)
			} else {
				c.podName = c.podmanContainer
			}
		case "kubectl":
			// We try the first item on purpose to make the command fail if not available
			podName, err := RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", c.kubernetesFilter, "-A", "-o=jsonpath={.items[0].metadata.name}")
			if err == nil {
				c.podName = string(podName[:])
			}
		}

	}

	return c.podName, err
}

// WaitForServer waits at most 60s for multi-user systemd target to be reached.
func (c *Connection) WaitForServer() {
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
		testCmd := exec.Command(command, args...)
		testCmd.Run()
		if testCmd.ProcessState.ExitCode() == 0 {
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatal().Msgf("Server didn't start within 60s")
}
