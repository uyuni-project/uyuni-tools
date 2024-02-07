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

type setCustomOptionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	Options               []string
}

func setCustomOptionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setCustomOptions",
		Short: "Set custom options for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setCustomOptionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setCustomOptions)
		},
	}

	cmd.Flags().String("KsLabel", "", "")
	cmd.Flags().String("Options", "", "$desc")

	return cmd
}

func setCustomOptions(globalFlags *types.GlobalFlags, flags *setCustomOptionsFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.Options)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
