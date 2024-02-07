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

type isErrataEmailNotifsForOrgFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func isErrataEmailNotifsForOrgCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isErrataEmailNotifsForOrg",
		Short: "Returns whether errata e-mail notifications are enabled
 for the organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isErrataEmailNotifsForOrgFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isErrataEmailNotifsForOrg)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func isErrataEmailNotifsForOrg(globalFlags *types.GlobalFlags, flags *isErrataEmailNotifsForOrgFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

