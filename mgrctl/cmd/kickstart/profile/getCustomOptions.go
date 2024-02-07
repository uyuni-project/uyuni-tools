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

type getCustomOptionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func getCustomOptionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getCustomOptions",
		Short: "Get custom options for a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getCustomOptionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getCustomOptions)
		},
	}

	cmd.Flags().String("KsLabel", "", "")

	return cmd
}

func getCustomOptions(globalFlags *types.GlobalFlags, flags *getCustomOptionsFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

