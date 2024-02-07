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

type listSystemsWithPackageFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Pid          int
	Name          string
	Version          string
	Release          string
}

func listSystemsWithPackageCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listSystemsWithPackage",
		Short: "Lists the systems that have the given installed package",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listSystemsWithPackageFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listSystemsWithPackage)
		},
	}

	cmd.Flags().String("Pid", "", "the package id")
	cmd.Flags().String("Name", "", "the package name")
	cmd.Flags().String("Version", "", "the package version")
	cmd.Flags().String("Release", "", "the package release")

	return cmd
}

func listSystemsWithPackage(globalFlags *types.GlobalFlags, flags *listSystemsWithPackageFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Pid, flags.Name, flags.Version, flags.Release)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

