package exec

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Envs        []string
	Interactive bool
	Tty         bool
}

// NewCommand returns a new cobra.Command for exec
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	execCmd := &cobra.Command{
		Use:   "exec '[command-to-run --with-args]'",
		Short: "execute commands inside the uyuni containers using 'sh -c'",
		Run: func(cmd *cobra.Command, args []string) {
			run(globalFlags, flags, cmd, args)
		},
	}
	execCmd.Flags().StringArrayVarP(&flags.Envs, "env", "e", []string{}, "environment variables to pass to the command")
	execCmd.Flags().BoolVarP(&flags.Interactive, "interactive", "i", false, "Pass stdin to the container")
	execCmd.Flags().BoolVarP(&flags.Tty, "tty", "t", false, "Stdin is a TTY")
	return execCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	command, podName := utils.GetPodName()

	commandArgs := []string{"exec"}
	if flags.Interactive {
		commandArgs = append(commandArgs, "-i")
	}
	if flags.Tty {
		commandArgs = append(commandArgs, "-t")
	}
	commandArgs = append(commandArgs, podName)

	if command == "kubectl" {
		commandArgs = append(commandArgs, "-c", "uyuni", "--")
	}

	newEnv := []string{}
	for _, env := range flags.Envs {
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
	if globalFlags.Verbose {
		fmt.Printf("> Running: %s %s\n", command, strings.Join(commandArgs, " "))
	}
	runCmd := exec.Command(command, commandArgs...)
	runCmd.Stdout = os.Stdout
	runCmd.Stdin = os.Stdin

	// Filter out kubectl line about terminated exit code
	stderr, err := runCmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err = runCmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "command terminated with exit code") {
			fmt.Fprintln(os.Stderr, line)
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	if err = runCmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			log.Fatal(err)
		}
	}
}
