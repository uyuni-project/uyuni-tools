package powermanagement

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/provisioning/powermanagement"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getStatusFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Name          string
}

func getStatusCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getStatus",
		Short: "Execute powermanagement actions",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getStatusFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getStatus)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Name", "", "")

	return cmd
}

func getStatus(globalFlags *types.GlobalFlags, flags *getStatusFlags, cmd *cobra.Command, args []string) error {

res, err := powermanagement.Powermanagement(&flags.ConnectionDetails, flags.Sid, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

