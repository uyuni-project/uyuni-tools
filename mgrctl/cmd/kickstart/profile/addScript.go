package profile

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addScriptFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel               string
	Name                  string
	Contents              string
	Interpreter           string
	Type                  string
	Chroot                bool
	Template              bool
	Erroronfail           bool
}

func addScriptCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addScript",
		Short: "Add a pre/post script to a kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addScriptFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addScript)
		},
	}

	cmd.Flags().String("KsLabel", "", "The kickstart label to add the script to.")
	cmd.Flags().String("Name", "", "The kickstart script name.")
	cmd.Flags().String("Contents", "", "The full script to add.")
	cmd.Flags().String("Interpreter", "", "The path to the interpreter to use (i.e. /bin/bash). An empty string will use the kickstart default interpreter.")
	cmd.Flags().String("Type", "", "The type of script (either 'pre' or 'post').")
	cmd.Flags().String("Chroot", "", "Whether to run the script in the chrooted install location (recommended) or not.")
	cmd.Flags().String("Template", "", "Enable templating using cobbler.")
	cmd.Flags().String("Erroronfail", "", "Whether to throw an error if the script fails or not")

	return cmd
}

func addScript(globalFlags *types.GlobalFlags, flags *addScriptFlags, cmd *cobra.Command, args []string) error {

	res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel, flags.Name, flags.Contents, flags.Interpreter, flags.Type, flags.Chroot, flags.Template, flags.Erroronfail)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
