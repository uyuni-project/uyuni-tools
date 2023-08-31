package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

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
			log.Fatalf("Backend command not found in PATH: %s\n", command)
		}
	case "":
		// Check kubectl with a timeout in case the configured cluster is not responding
		_, err := exec.LookPath("kubectl")
		if err == nil {
			if err = exec.Command("kubectl", "--request-timeout=30s", "get", "pod").Run(); err != nil {
				log.Print("kubectl not configured to connect to a cluster, ignoring")
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

		log.Fatal("Neither podman, podman-remote nor kubectl is available")
	default:
		log.Fatalf("Unsupported backend %s\n", backend)
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
				log.Fatalf("Container %s is not running on podman", pod)
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
	log.Fatalf("Server didn't start within 60s")
}

func RunCmd(command string, args []string, errMessage string, verbose bool) {
	if verbose {
		fmt.Printf("> Running: %s %s\n", command, strings.Join(args, " "))
	}
	cmd := exec.Command(command, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("%s:\n  %s\n", errMessage, strings.ReplaceAll(string(out[:]), "\n", "\n  "))
	}
}

const PROMPT_END = ": "

func AskPasswordIfMissing(value *string, prompt string) {
	if *value == "" {
		fmt.Print(prompt + PROMPT_END)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Failed to read password: %s\n", err)
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
			log.Fatalf("Failed to read input: %s\n", err)
		}
		*value = newValue
		fmt.Println()
	}
}

// Get the timezone set on the machine running the tool
func GetLocalTimezone() string {
	out, err := exec.Command("timedatectl", "show", "--value", "-p", "Timezone").Output()
	if err != nil {
		log.Fatalf("Failed to run timedatectl show --value -p Timezone: %s\n", err)
	}
	return string(out)
}

// Check if a given path exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else if !os.IsNotExist(err) {
		log.Fatalf("Failed to stat %s file: %s\n", path, err)
	}
	return false
}
