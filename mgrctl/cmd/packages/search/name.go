package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages/search"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type nameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func nameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "name",
		Short: "Search the lucene package indexes for all packages which
          match the given name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags nameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, name)
		},
	}

	cmd.Flags().String("Name", "", "package name to search for")

	return cmd
}

func name(globalFlags *types.GlobalFlags, flags *nameFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

