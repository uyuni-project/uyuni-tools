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

type getPolicyForScapFileUploadFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func getPolicyForScapFileUploadCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPolicyForScapFileUpload",
		Short: "Get the status of SCAP detailed result file upload settings
 for the given organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPolicyForScapFileUploadFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPolicyForScapFileUpload)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func getPolicyForScapFileUpload(globalFlags *types.GlobalFlags, flags *getPolicyForScapFileUploadFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

