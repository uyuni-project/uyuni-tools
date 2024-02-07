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

type searchByNameFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Regexp          string
}

func searchByNameCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "searchByName",
		Short: "Returns a list of system IDs whose name matches
  the supplied regular expression(defined by
  http://docs.oracle.com/javase/1.5.0/docs/api/java/util/regex/Pattern.html_blank
 Java representation of regular expressions)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags searchByNameFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, searchByName)
		},
	}

	cmd.Flags().String("Regexp", "", "A regular expression")

	return cmd
}

func searchByName(globalFlags *types.GlobalFlags, flags *searchByNameFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Regexp)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

