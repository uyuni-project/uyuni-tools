package api

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyunictl/shared/utils"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	apiFlags := &api.ConnectionDetails{}

	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "JSON over HTTP API helper tool",
	}

	apiGet := &cobra.Command{
		Use:   "get path [parameters]...",
		Short: "Call API GET request",
		Long:  "Takes an API path and optional parameters and then issues GET request with the specified path and parameters. If user and password are provided, calls login before API call",
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "ctlconfig", cmd)
			if err := viper.Unmarshal(&apiFlags); err != nil {
				log.Fatal().Err(err).Msgf("Failed to unmarshall configuration")
			}
			runGet(globalFlags, apiFlags, cmd, args)
		},
	}

	apiPost := &cobra.Command{
		Use:   "post path parameters...",
		Short: "Call API POST request",
		Long:  "Takes an API path and parameters and then issues POST request with the specified path and parameters. User and password are mandatory. Parameters can be either JSON encoded string or one or more key=value pairs.",
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "ctlconfig", cmd)
			if err := viper.Unmarshal(&apiFlags); err != nil {
				log.Fatal().Err(err).Msgf("Failed to unmarshall configuration")
			}
			runPost(globalFlags, apiFlags, cmd, args)
		},
	}

	apiCmd.AddCommand(apiGet)
	apiCmd.AddCommand(apiPost)

	cmd_utils.AddAPIFlags(apiCmd, apiFlags)
	return apiCmd
}
