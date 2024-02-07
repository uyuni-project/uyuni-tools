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

type setPolicyForScapFileUploadFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func setPolicyForScapFileUploadCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPolicyForScapFileUpload",
		Short: "Set the status of SCAP detailed result file upload settings
 for the given organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setPolicyForScapFileUploadFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setPolicyForScapFileUpload)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func setPolicyForScapFileUpload(globalFlags *types.GlobalFlags, flags *setPolicyForScapFileUploadFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

