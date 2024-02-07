package errata

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/errata"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deleteFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	AdvisoryName          string
}

func deleteCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an erratum.  This method will only allow for deletion
 of custom errata created either through the UI or API.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, delete)
		},
	}

	cmd.Flags().String("AdvisoryName", "", "")

	return cmd
}

func delete(globalFlags *types.GlobalFlags, flags *deleteFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.AdvisoryName)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

