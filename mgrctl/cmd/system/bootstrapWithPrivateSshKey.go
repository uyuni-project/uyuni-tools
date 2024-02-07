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

type bootstrapWithPrivateSshKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Host          string
	SshPort          int
	SshUser          string
	SshPrivKey          string
	SshPrivKeyPass          string
	ActivationKey          string
	SaltSSH          bool
	ProxyId          int
	ReactivationKey          string
}

func bootstrapWithPrivateSshKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrapWithPrivateSshKey",
		Short: "Bootstrap a system for management via either Salt or Salt SSH.
 Use SSH private key for authentication.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags bootstrapWithPrivateSshKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, bootstrapWithPrivateSshKey)
		},
	}

	cmd.Flags().String("Host", "", "Hostname or IP address of target")
	cmd.Flags().String("SshPort", "", "SSH port on target machine")
	cmd.Flags().String("SshUser", "", "SSH user on target machine")
	cmd.Flags().String("SshPrivKey", "", "SSH private key as a string in PEM format")
	cmd.Flags().String("SshPrivKeyPass", "", "SSH passphrase for the key (use empty string for no passphrase)")
	cmd.Flags().String("ActivationKey", "", "Activation key")
	cmd.Flags().String("SaltSSH", "", "Manage system with Salt SSH")
	cmd.Flags().String("ProxyId", "", "System ID of proxy to use")
	cmd.Flags().String("ReactivationKey", "", "Reactivation key")

	return cmd
}

func bootstrapWithPrivateSshKey(globalFlags *types.GlobalFlags, flags *bootstrapWithPrivateSshKeyFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Host, flags.SshPort, flags.SshUser, flags.SshPrivKey, flags.SshPrivKeyPass, flags.ActivationKey, flags.SaltSSH, flags.ProxyId, flags.ReactivationKey)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

