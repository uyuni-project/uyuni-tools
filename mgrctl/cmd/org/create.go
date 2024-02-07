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

type createFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgName          string
	AdminLogin          string
	AdminPassword          string
	Prefix          string
	FirstName          string
	LastName          string
	Email          string
	UsePamAuth          bool
}

func createCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new organization and associated administrator account.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, create)
		},
	}

	cmd.Flags().String("OrgName", "", "Organization name. Must meet same criteria as in the web UI.")
	cmd.Flags().String("AdminLogin", "", "New administrator login name.")
	cmd.Flags().String("AdminPassword", "", "New administrator password.")
	cmd.Flags().String("Prefix", "", "New administrator's prefix. Must match one of the values available in the web UI. (i.e. Dr., Mr., Mrs., Sr., etc.)")
	cmd.Flags().String("FirstName", "", "New administrator's first name.")
	cmd.Flags().String("LastName", "", "New administrator's first name.")
	cmd.Flags().String("Email", "", "New administrator's e-mail.")
	cmd.Flags().String("UsePamAuth", "", "True if PAM authentication should be used for the new administrator account.")

	return cmd
}

func create(globalFlags *types.GlobalFlags, flags *createFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgName, flags.AdminLogin, flags.AdminPassword, flags.Prefix, flags.FirstName, flags.LastName, flags.Email, flags.UsePamAuth)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

