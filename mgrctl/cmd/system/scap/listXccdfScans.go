package scap

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/scap"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listXccdfScansFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
}

func listXccdfScansCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listXccdfScans",
		Short: "Return a list of finished OpenSCAP scans for a given system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listXccdfScansFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listXccdfScans)
		},
	}

	cmd.Flags().String("Sid", "", "")

	return cmd
}

func listXccdfScans(globalFlags *types.GlobalFlags, flags *listXccdfScansFlags, cmd *cobra.Command, args []string) error {

res, err := scap.Scap(&flags.ConnectionDetails, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

