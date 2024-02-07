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

type getRelevantErrataByTypeFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	AdvisoryType          string
}

func getRelevantErrataByTypeCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getRelevantErrataByType",
		Short: "Returns a list of all errata of the specified type that are
 relevant to the system.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getRelevantErrataByTypeFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getRelevantErrataByType)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("AdvisoryType", "", "type of advisory (one of of the following: 'Security Advisory', 'Product Enhancement Advisory', 'Bug Fix Advisory'")

	return cmd
}

func getRelevantErrataByType(globalFlags *types.GlobalFlags, flags *getRelevantErrataByTypeFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.AdvisoryType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

