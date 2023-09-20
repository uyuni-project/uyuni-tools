package exec

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyunictl/shared/utils"
)

type flagpole struct {
	Envs        []string `mapstructure:"env"`
	Interactive bool
	Tty         bool
	Backend     string
}

// NewCommand returns a new cobra.Command for exec
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	execCmd := &cobra.Command{
		Use:   "exec '[command-to-run --with-args]'",
		Short: "Execute commands inside the uyuni containers using 'sh -c'",
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "ctlconfig", cmd)
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msgf("Failed to unmarshall configuration")
			}
			run(globalFlags, flags, cmd, args)
		},
	}
	execCmd.Flags().StringSliceP("env", "e", []string{}, "environment variables to pass to the command, separated by commas")
	execCmd.Flags().BoolP("interactive", "i", false, "Pass stdin to the container")
	execCmd.Flags().BoolP("tty", "t", false, "Stdin is a TTY")

	cmd_utils.AddBackendFlag(execCmd)
	return execCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {

	command, podName := utils.GetPodName(globalFlags, flags.Backend, true)

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
	for _, envValue := range flags.Envs {
		if !strings.Contains(envValue, "=") {
			if value, set := os.LookupEnv(envValue); set {
				newEnv = append(newEnv, fmt.Sprintf("%s=%s", envValue, value))
			}
		} else {
			newEnv = append(newEnv, envValue)
		}
	}
	if len(newEnv) > 0 {
		commandArgs = append(commandArgs, "env")
		commandArgs = append(commandArgs, newEnv...)
	}
	commandArgs = append(commandArgs, "sh", "-c", strings.Join(args, " "))
	err := RunRawCmd(command, commandArgs, false)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
	}
}

func RunRawCmd(command string, args []string, outputToLog bool) error {

	log.Debug().Msgf(" Running: %s %s", command, strings.Join(args, " "))

	runCmd := exec.Command(command, args...)
	runCmd.Stdin = os.Stdin
	if outputToLog {
		runCmd.Stdout = utils.OutputLogWriter{Logger: log.Logger, LogLevel: zerolog.DebugLevel}
	} else {
		runCmd.Stdout = os.Stdout
	}

	// Filter out kubectl line about terminated exit code
	stderr, err := runCmd.StderrPipe()
	if err != nil {
		log.Debug().Err(err).Msg("error starting stderr processor for command")
		return err
	}
	defer stderr.Close()

	if err = runCmd.Start(); err != nil {
		log.Debug().Err(err).Msg("error starting command")
		return err
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		// needed because of kubernetes installation, to ignore the stderr
		if !strings.HasPrefix(line, "command terminated with exit code") {
			if outputToLog {
				log.Debug().Msg(line)
			} else {
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}

	if scanner.Err() != nil {
		log.Debug().Msg("error scanning stderr")
		return scanner.Err()
	}
	return runCmd.Wait()
}
