package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func GetPodName() (string, string) {
	command := "podman"
	pod := "uyuni-server"

	_, err := exec.LookPath("kubectl")
	if err == nil {
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
