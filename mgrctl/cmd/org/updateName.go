package org

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/org"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateNameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId                 int
	Name                  string
}

func updateNameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateName",
		Short: "Updates the name of an organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateNameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateName)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("Name", "", "Organization name. Must meet same criteria as in the web UI.")

	return cmd
}

func updateName(globalFlags *types.GlobalFlags, flags *updateNameFlags, cmd *cobra.Command, args []string) error {

	res, err := org.Org(&flags.ConnectionDetails, flags.OrgId, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
