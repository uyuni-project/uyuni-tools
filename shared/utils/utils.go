package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
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

func RunCmd(command string, args []string, errMessage string, verbose bool) {
	if verbose {
		fmt.Printf("> Running: %s %s\n", command, strings.Join(args, " "))
	}
	cmd := exec.Command(command, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalf("%s:\n  %s\n", errMessage, strings.ReplaceAll(string(out[:]), "\n", "\n  "))
	}
}
