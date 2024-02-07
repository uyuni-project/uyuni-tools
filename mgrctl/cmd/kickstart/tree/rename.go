package tree

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/tree"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type renameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OriginalLabel         string
	NewLabel              string
}

func renameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename",
		Short: "Rename a Kickstart Tree (Distribution) in #product().",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags renameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, rename)
		},
	}

	cmd.Flags().String("OriginalLabel", "", "Label for the kickstart tree to rename.")
	cmd.Flags().String("NewLabel", "", "The kickstart tree's new label.")

	return cmd
}

func rename(globalFlags *types.GlobalFlags, flags *renameFlags, cmd *cobra.Command, args []string) error {

	res, err := tree.Tree(&flags.ConnectionDetails, flags.OriginalLabel, flags.NewLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
