package scap

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/scap"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type scheduleXccdfScanFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids          []int
	XccdfPath          string
	OscapParams          string
	Date          $date
	XccdfPath          string
	OscapPrams          string
	OvalFiles          string
	Sid          int
}

func scheduleXccdfScanCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduleXccdfScan",
		Short: "Schedule OpenSCAP scan.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags scheduleXccdfScanFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, scheduleXccdfScan)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("XccdfPath", "", "path to xccdf content on targeted systems.")
	cmd.Flags().String("OscapParams", "", "additional parameters for oscap tool.")
	cmd.Flags().String("Date", "", "The date to schedule the action")
	cmd.Flags().String("XccdfPath", "", "Path to xccdf content on targeted systems.")
	cmd.Flags().String("OscapPrams", "", "Additional parameters for oscap tool.")
	cmd.Flags().String("OvalFiles", "", "Additional OVAL files for oscap tool.")
	cmd.Flags().String("Sid", "", "")

	return cmd
}

func scheduleXccdfScan(globalFlags *types.GlobalFlags, flags *scheduleXccdfScanFlags, cmd *cobra.Command, args []string) error {

res, err := scap.Scap(&flags.ConnectionDetails, flags.Sids, flags.XccdfPath, flags.OscapParams, flags.Date, flags.XccdfPath, flags.OscapPrams, flags.OvalFiles, flags.Sid)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

