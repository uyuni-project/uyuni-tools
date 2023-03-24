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
			commandArgs = []string{"exec", podName}
		case "kubectl":
			commandArgs = []string{"exec", podName, "-c", "uyuni", "--"}
		default:
			log.Fatalf("Unknown container kind: %s", command)
		}

		newEnv := []string{}
		for _, env := range envs {
			if !strings.Contains(env, "=") {
				if value, set := os.LookupEnv(env); set {
					newEnv = append(newEnv, fmt.Sprintf("%s=%s", env, value))
				}
			} else {
				newEnv = append(newEnv, env)
			}
		}
		if len(newEnv) > 0 {
			commandArgs = append(commandArgs, "env")
			commandArgs = append(commandArgs, newEnv...)
		}
		commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))
		if Verbose {
			fmt.Printf("> Running: %s %s\n", command, strings.Join(commandArgs, " "))
		}
		runCmd := exec.Command(command, commandArgs...)
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

var envs []string

func init() {
	execCmd.Flags().StringArrayVarP(&envs, "env", "e", []string{}, "environment variables to pass to the command")
	rootCmd.AddCommand(execCmd)
}
