package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getRegistrationDateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid                   int
}

func getRegistrationDateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRegistrationDate",
		Short: "Returns the date the system was registered.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRegistrationDateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRegistrationDate)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func getRegistrationDate(globalFlags *types.GlobalFlags, flags *getRegistrationDateFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
