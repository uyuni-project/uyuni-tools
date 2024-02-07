package proxy

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/proxy"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type containerConfigFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProxyName          string
	ProxyPort          int
	Server          string
	MaxCache          int
	Email          string
	RootCA          string
	ProxyCrt          string
	ProxyKey          string
	CaCrt          string
	CaKey          string
	CaPassword          string
	$param.getFlagName()          $param.getType()
	Country          string
	State          string
	City          string
	Org          string
	OrgUnit          string
	SslEmail          string
}

func containerConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "containerConfig",
		Short: "Compute and download the configuration for proxy containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags containerConfigFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, containerConfig)
		},
	}

	cmd.Flags().String("ProxyName", "", "")
	cmd.Flags().String("ProxyPort", "", "")
	cmd.Flags().String("Server", "", "")
	cmd.Flags().String("MaxCache", "", "")
	cmd.Flags().String("Email", "", "")
	cmd.Flags().String("RootCA", "", "")
	cmd.Flags().String("ProxyCrt", "", "")
	cmd.Flags().String("ProxyKey", "", "")
	cmd.Flags().String("CaCrt", "", "")
	cmd.Flags().String("CaKey", "", "")
	cmd.Flags().String("CaPassword", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("Country", "", "")
	cmd.Flags().String("State", "", "")
	cmd.Flags().String("City", "", "")
	cmd.Flags().String("Org", "", "")
	cmd.Flags().String("OrgUnit", "", "")
	cmd.Flags().String("SslEmail", "", "")

	return cmd
}

func containerConfig(globalFlags *types.GlobalFlags, flags *containerConfigFlags, cmd *cobra.Command, args []string) error {

res, err := proxy.Proxy(&flags.ConnectionDetails, flags.ProxyName, flags.ProxyPort, flags.Server, flags.MaxCache, flags.Email, flags.RootCA, flags.ProxyCrt, flags.ProxyKey, flags.CaCrt, flags.CaKey, flags.CaPassword, flags.$param.getFlagName(), flags.Country, flags.State, flags.City, flags.Org, flags.OrgUnit, flags.SslEmail)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

