package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setLoggingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	Pre          bool
	Post          bool
}

func setLoggingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setLogging",
		Short: "Set logging options for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setLoggingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setLogging)
		},
	}

	cmd.Flags().String("KsLabel", "", "Label of kickstart profile to be changed.")
	cmd.Flags().String("Pre", "", "whether or not to log      the pre section of a kickstart to /root/ks-pre.log")
	cmd.Flags().String("Post", "", "whether or not to log      the post section of a kickstart to /root/ks-post.log")

	return cmd
}

func setLogging(globalFlags *types.GlobalFlags, flags *setLoggingFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.Pre, flags.Post)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

