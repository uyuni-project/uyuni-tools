package cp

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	cmd_utils "github.com/uyuni-project/uyuni-tools/uyunictl/shared/utils"
)

type flagpole struct {
	User    string
	Group   string
	Backend string
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	cpCmd := &cobra.Command{
		Use:   "cp [path/to/source.file] [path/to/destination.file]",
		Short: "copy files to and from the containers",
		Long: `Takes a source and destination parameters.
	One of them can be prefixed with 'server:' to indicate the path is within the server pod.`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			viper := utils.ReadConfig(globalFlags.ConfigPath, "ctlconfig", cmd)
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msgf("Failed to unmarshall configuration")
			}
			run(flags, cmd, args)
		},
	}

	cpCmd.Flags().String("user", "", "User or UID to set on the destination file")
	cpCmd.Flags().String("group", "susemanager", "Group or GID to set on the destination file")

	cmd_utils.AddBackendFlag(cpCmd)
	return cpCmd
}

func run(flags *flagpole, cmd *cobra.Command, args []string) {
	cnx := utils.NewConnection(flags.Backend)
	utils.Copy(cnx, args[0], args[1], flags.User, flags.Group)
}
