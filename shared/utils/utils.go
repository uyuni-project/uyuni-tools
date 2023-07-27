package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/term"
)

func GetCommand() string {
	command := ""

	_, err := exec.LookPath("kubectl")
	if err == nil {
		if err = exec.Command("kubectl", "get", "pod").Run(); err != nil {
			log.Print("kubectl not configured to connect to a cluster, ignoring")
		} else {
			command = "kubectl"
		}
	}

	if _, err = exec.LookPath("podman"); err == nil {
		command = "podman"
	}

	if command == "" {
		log.Fatal("Neither podman nor kubectl are available")
	}
	return command
}

func GetPodName() (string, string) {
	command := GetCommand()
	pod := "uyuni-server"

	switch command {
	case "podman":
		if out, _ := exec.Command("podman", "ps", "-q", "-f", "name="+pod).Output(); len(out) == 0 {
			log.Fatalf("Container %s is not running on podman", pod)
		}
	case "kubectl":
		podCmd := exec.Command("kubectl", "get", "pod", "-lapp=uyuni", "-o=jsonpath={.items[0].metadata.name}")
		podName, err := podCmd.Output()
		if err == nil {
			command = "kubectl"
			pod = string(podName[:])
		}
	}
	return command, pod
}

// WaitForServer waits at most 60s for multi-user systemd target to be reached.
func WaitForServer() {
	cmd, podName := GetPodName()
	// Wait for the system to be up
	for i := 0; i < 60; i++ {
		cmd := exec.Command(cmd, "exec", podName, "--", "systemctl", "is-active", "-q", "multi-user.target")
		cmd.Run()
		if cmd.ProcessState.ExitCode() == 0 {
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

func AskPasswordIfMissing(viper *viper.Viper, key string, prompt string) {
	value := viper.GetString(key)
	if value == "" {
		fmt.Print(prompt)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("Failed to read password: %s\n", err)
		}
		viper.Set(key, string(bytePassword))
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
