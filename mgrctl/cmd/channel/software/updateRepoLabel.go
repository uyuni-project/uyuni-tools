package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateRepoLabelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id                    int
	Label                 string
	Label                 string
	NewLabel              string
}

func updateRepoLabelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateRepoLabel",
		Short: "Updates repository label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateRepoLabelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateRepoLabel)
		},
	}

	cmd.Flags().String("Id", "", "repository ID")
	cmd.Flags().String("Label", "", "new repository label")
	cmd.Flags().String("Label", "", "repository label")
	cmd.Flags().String("NewLabel", "", "new repository label")

	return cmd
}

func updateRepoLabel(globalFlags *types.GlobalFlags, flags *updateRepoLabelFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.Id, flags.Label, flags.Label, flags.NewLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
