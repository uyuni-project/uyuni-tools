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

type setErrataEmailNotifsForOrgFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
	Enable          bool
}

func setErrataEmailNotifsForOrgCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setErrataEmailNotifsForOrg",
		Short: "Dis/enables errata e-mail notifications for the organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setErrataEmailNotifsForOrgFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setErrataEmailNotifsForOrg)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("Enable", "", "Use true/false to enable/disable")

	return cmd
}

func setErrataEmailNotifsForOrg(globalFlags *types.GlobalFlags, flags *setErrataEmailNotifsForOrgFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId, flags.Enable)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

