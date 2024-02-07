package packages

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type findByNvreaFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	Version          string
	Release          string
	Epoch          string
	ArchLabel          string
}

func findByNvreaCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "findByNvrea",
		Short: "Lookup the details for packages with the given name, version,
          release, architecture label, and (optionally) epoch.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags findByNvreaFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, findByNvrea)
		},
	}

	cmd.Flags().String("Name", "", "")
	cmd.Flags().String("Version", "", "")
	cmd.Flags().String("Release", "", "")
	cmd.Flags().String("Epoch", "", "If set to something other than empty string,          strict matching will be used and the epoch string must be correct.          If set to an empty string, if the epoch is null or there is only one          NVRA combination, it will be returned.  (Empty string is recommended.)")
	cmd.Flags().String("ArchLabel", "", "")

	return cmd
}

func findByNvrea(globalFlags *types.GlobalFlags, flags *findByNvreaFlags, cmd *cobra.Command, args []string) error {

res, err := packages.Packages(&flags.ConnectionDetails, flags.Name, flags.Version, flags.Release, flags.Epoch, flags.ArchLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

