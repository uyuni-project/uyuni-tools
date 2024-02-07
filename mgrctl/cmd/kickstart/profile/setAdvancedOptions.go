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

type setAdvancedOptionsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func setAdvancedOptionsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setAdvancedOptions",
		Short: "Set advanced options for a kickstart profile.
 'md5_crypt_rootpw' is not supported anymore.
 If 'sha256_crypt_rootpw' is set to 'True', 'root_pw' is taken as plaintext and
 will sha256 encrypted on server side, otherwise a hash encoded password
 (according to the auth option) is expected",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setAdvancedOptionsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setAdvancedOptions)
		},
	}

	cmd.Flags().String("KsLabel", "", "")

	return cmd
}

func setAdvancedOptions(globalFlags *types.GlobalFlags, flags *setAdvancedOptionsFlags, cmd *cobra.Command, args []string) error {

res, err := profile.Profile(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

