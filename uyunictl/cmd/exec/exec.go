package exec

import (
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
	err := utils.Exec(globalFlags, flags.Backend, flags.Interactive, flags.Tty, false, flags.Envs, args...)
	if err != nil {
		log.Debug().Err(err).Msg("error running the command")
	}
}
