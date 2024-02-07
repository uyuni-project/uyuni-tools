package delta

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/delta"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SourceImageId         int
	TargetImageId         int
}

func getDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getDetails",
		Short: "Get details of an Image",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getDetails)
		},
	}

	cmd.Flags().String("SourceImageId", "", "")
	cmd.Flags().String("TargetImageId", "", "")

	return cmd
}

func getDetails(globalFlags *types.GlobalFlags, flags *getDetailsFlags, cmd *cobra.Command, args []string) error {

	res, err := delta.Delta(&flags.ConnectionDetails, flags.SourceImageId, flags.TargetImageId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
