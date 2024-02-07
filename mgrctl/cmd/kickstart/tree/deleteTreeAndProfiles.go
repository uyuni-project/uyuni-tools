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

type deleteTreeAndProfilesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	TreeLabel          string
}

func deleteTreeAndProfilesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteTreeAndProfiles",
		Short: "Delete a kickstarttree and any profiles associated with
 this kickstart tree.  WARNING:  This will delete all profiles
 associated with this kickstart tree!",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteTreeAndProfilesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteTreeAndProfiles)
		},
	}

	cmd.Flags().String("TreeLabel", "", "Label for the kickstart tree to delete.")

	return cmd
}

func deleteTreeAndProfiles(globalFlags *types.GlobalFlags, flags *deleteTreeAndProfilesFlags, cmd *cobra.Command, args []string) error {

res, err := tree.Tree(&flags.ConnectionDetails, flags.TreeLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

