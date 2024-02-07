package slave

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/slave"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SlaveFqdn             string
	IsEnabled             bool
	AllowAllOrgs          bool
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Slave, known to this Master.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("SlaveFqdn", "", "Slave's fully-qualified domain name")
	cmd.Flags().String("IsEnabled", "", "Let this slave talk to us?")
	cmd.Flags().String("AllowAllOrgs", "", "Export all our orgs to this slave?")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

	res, err := slave.Slave(&flags.ConnectionDetails, flags.SlaveFqdn, flags.IsEnabled, flags.AllowAllOrgs)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
