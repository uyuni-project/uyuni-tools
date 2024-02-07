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

type deleteProfileFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
}

func deleteProfileCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteProfile",
		Short: "Delete a kickstart profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteProfileFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteProfile)
		},
	}

	cmd.Flags().String("KsLabel", "", "The label of the kickstart profile you want to remove")

	return cmd
}

func deleteProfile(globalFlags *types.GlobalFlags, flags *deleteProfileFlags, cmd *cobra.Command, args []string) error {

	res, err := kickstart.Kickstart(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
