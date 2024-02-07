package recurring

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/recurring"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listByEntityFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	$param.getFlagName()          $param.getType()
	Id          int
}

func listByEntityCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listByEntity",
		Short: "Return a list of recurring actions for a given entity.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listByEntityFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listByEntity)
		},
	}

	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("Id", "", "the ID of the target entity")

	return cmd
}

func listByEntity(globalFlags *types.GlobalFlags, flags *listByEntityFlags, cmd *cobra.Command, args []string) error {

res, err := recurring.Recurring(&flags.ConnectionDetails, flags.$param.getFlagName(), flags.Id)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

