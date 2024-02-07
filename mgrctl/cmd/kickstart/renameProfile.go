package kickstart

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type renameProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OriginalLabel          string
	NewLabel          string
}

func renameProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renameProfile",
		Short: "Rename a kickstart profile in #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags renameProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, renameProfile)
		},
	}

	cmd.Flags().String("OriginalLabel", "", "Label for the kickstart profile you want to rename")
	cmd.Flags().String("NewLabel", "", "new label to change to")

	return cmd
}

func renameProfile(globalFlags *types.GlobalFlags, flags *renameProfileFlags, cmd *cobra.Command, args []string) error {

res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.OriginalLabel, flags.NewLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

