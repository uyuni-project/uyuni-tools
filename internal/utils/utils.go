package utils

import "os/exec"

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
