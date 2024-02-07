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

type updateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SlaveId          int
	SlaveFqdn          string
	IsEnabled          bool
	AllowAllOrgs          bool
}

func updateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates attributes of the specified Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, update)
		},
	}

	cmd.Flags().String("SlaveId", "", "ID of the Slave to update")
	cmd.Flags().String("SlaveFqdn", "", "Slave's fully-qualified domain name")
	cmd.Flags().String("IsEnabled", "", "Let this slave talk to us?")
	cmd.Flags().String("AllowAllOrgs", "", "Export all our orgs to this Slave?")

	return cmd
}

func update(globalFlags *types.GlobalFlags, flags *updateFlags, cmd *cobra.Command, args []string) error {

res, err := slave.Slave(&flags.ConnectionDetails, flags.SlaveId, flags.SlaveFqdn, flags.IsEnabled, flags.AllowAllOrgs)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

