package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type alignMetadataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelFromLabel      string
	ChannelToLabel        string
	MetadataType          string
}

func alignMetadataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alignMetadata",
		Short: "Align the metadata of a channel to another channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags alignMetadataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, alignMetadata)
		},
	}

	cmd.Flags().String("ChannelFromLabel", "", "the label of the source channel")
	cmd.Flags().String("ChannelToLabel", "", "the label of the target channel")
	cmd.Flags().String("MetadataType", "", "the metadata type. Only 'modules' supported currently.")

	return cmd
}

func alignMetadata(globalFlags *types.GlobalFlags, flags *alignMetadataFlags, cmd *cobra.Command, args []string) error {

	res, err := software.Software(&flags.ConnectionDetails, flags.ChannelFromLabel, flags.ChannelToLabel, flags.MetadataType)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
