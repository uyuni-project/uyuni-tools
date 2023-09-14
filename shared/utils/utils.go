package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"golang.org/x/term"
)

func GetCommand(backend string) string {
	command := ""

	switch backend {
	case "podman":
		fallthrough
	case "podman-remote":
		fallthrough
	case "kubectl":
		command = backend
		if _, err := exec.LookPath(command); err != nil {
			log.Fatal().Msgf("Backend command not found in PATH: %s", command)
		}
	case "":
		// Check kubectl with a timeout in case the configured cluster is not responding
		_, err := exec.LookPath("kubectl")
		if err == nil {
			if err = exec.Command("kubectl", "--request-timeout=30s", "get", "pod").Run(); err != nil {
				log.Info().Msg("kubectl not configured to connect to a cluster, ignoring")
			} else {
				return "kubectl"
			}
		}

		// Search for other backends
		bins := []string{"podman", "podman-remote"}
		for _, bin := range bins {
			if _, err := exec.LookPath(bin); err == nil {
				return bin
			}
		}

		log.Fatal().Msg("Neither podman, podman-remote nor kubectl is available")
	default:
		log.Fatal().Msgf("Unsupported backend %s", backend)
	}
	return command
}

func GetPodName(globalFlags *types.GlobalFlags, backend string, fail bool) (string, string) {
	command := GetCommand(backend)
	pod := "uyuni-server"

	switch command {
	case "podman-remote":
		fallthrough
	case "podman":
		if out, _ := exec.Command(command, "ps", "-q", "-f", "name="+pod).Output(); len(out) == 0 {
			if fail {
				log.Fatal().Msgf("Container %s is not running on podman", pod)
			}
		}
	case "kubectl":
		podCmd := exec.Command("kubectl", "get", "pod", "-lapp=uyuni", "-o=jsonpath={.items[0].metadata.name}")
		podName, err := podCmd.Output()
		if err == nil {
			pod = string(podName[:])
		}
	}
	return command, pod
}

// WaitForServer waits at most 60s for multi-user systemd target to be reached.
func WaitForServer(globalFlags *types.GlobalFlags, backend string) {
	// Wait for the system to be up
	for i := 0; i < 60; i++ {
		cmd, podName := GetPodName(globalFlags, backend, false)
		args := []string{"exec", podName}
		if cmd == "kubectl" {
			args = append(args, "--")
		}
		args = append(args, "systemctl", "is-active", "-q", "multi-user.target")
		testCmd := exec.Command(cmd, args...)
		testCmd.Run()
		if testCmd.ProcessState.ExitCode() == 0 {
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatal().Msgf("Server didn't start within 60s")
}

const PROMPT_END = ": "

func AskPasswordIfMissing(value *string, prompt string) {
	if *value == "" {
		fmt.Print(prompt + PROMPT_END)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read password")
		}
		*value = string(bytePassword)
		fmt.Println()
	}
}

func AskIfMissing(value *string, prompt string) {
	if *value == "" {
		fmt.Print(prompt + PROMPT_END)
		reader := bufio.NewReader(os.Stdin)
		newValue, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal().Err(err).Msgf("Failed to read input")
		}
		*value = newValue
		fmt.Println()
	}
}

// Get the timezone set on the machine running the tool
func GetLocalTimezone() string {
	out, err := exec.Command("timedatectl", "show", "--value", "-p", "Timezone").Output()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to run timedatectl show --value -p Timezone")
	}
	return string(out)
}

// Check if a given path exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if !os.IsNotExist(err) {
		log.Fatal().Err(err).Msgf("Failed to stat %s file", path)
	}
	return false
}
