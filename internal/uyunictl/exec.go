package uyunictl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/internal/utils"
)

var execCmd = &cobra.Command{
	Use:   "exec '[command-to-run --with-args]'",
	Short: "execute commands inside the uyuni containers using 'sh -c'",
	Run: func(cmd *cobra.Command, args []string) {
		command, podName := utils.GetPodName()
		commandArgs := []string{}

		switch command {
		case "podman":
			commandArgs = []string{"exec", podName, "sh", "-c"}
		case "kubectl":
			commandArgs = []string{"exec", podName, "-c", "uyuni", "--", "sh", "-c"}
		default:
			log.Fatalf("Unknown container kind: %s", command)
		}

		allArgs := append(commandArgs, args...)
		if Verbose {
			fmt.Printf("> Running: %s %s\n", command, strings.Join(allArgs, " "))
		}
		runCmd := exec.Command(command, allArgs...)
		runCmd.Stderr = os.Stderr
		runCmd.Stdout = os.Stdout
		err := runCmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ExitCode())
			} else {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
