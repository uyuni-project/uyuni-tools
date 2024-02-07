package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getEncodedFileRevisionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	FilePath          string
	Revision          int
}

func getEncodedFileRevisionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getEncodedFileRevision",
		Short: "Get revision of the specified configuration file and transmit the
             contents as base64 encoded.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getEncodedFileRevisionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getEncodedFileRevision)
		},
	}

	cmd.Flags().String("Label", "", "label of config channel to lookup on")
	cmd.Flags().String("FilePath", "", "config file path to examine")
	cmd.Flags().String("Revision", "", "config file revision to examine")

	return cmd
}

func getEncodedFileRevision(globalFlags *types.GlobalFlags, flags *getEncodedFileRevisionFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.FilePath, flags.Revision)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

