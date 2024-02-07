package contentmanagement

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/contentmanagement"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
	Rule          string
	EntityType          string
	$param.getFlagName()          $param.getType()
}

func createFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createFilter",
		Short: "Create a Content Filter
 #paragraph_end()
 #paragraph()
 The following filters are available (you can get the list in machine-readable format using
 the listFilterCriteria() endpoint):
 #paragraph_end()
 #paragraph()
 Package filtering:
 #itemlist()
   #item("by name - field: name; matchers: contains or matches")
   #item("by name, epoch, version, release and architecture - field: nevr or nevra; matcher: equals")
  #itemlist_end()
 #paragraph_end()
 #paragraph()
 Errata/Patch filtering:
 #itemlist()
   #item("by advisory name - field: advisory_name; matcher: equals or matches")
   #item("by type - field: advisory_type (e.g. 'Security Advisory'); matcher: equals")
   #item("by synopsis - field: synopsis; matcher: equals, contains or matches")
   #item("by keyword - field: keyword; matcher: contains")
   #item("by date - field: issue_date; matcher: greater or greatereq; value needs to be in ISO format e.g
   2022-12-10T12:00:00Z")
   #item("by affected package name - field: package_name; matcher: contains_pkg_name or matches_pkg_name")
   #item("by affected package with version - field: package_nevr; matcher: contains_pkg_lt_evr,
   contains_pkg_le_evr, contains_pkg_eq_evr, contains_pkg_ge_evr or contains_pkg_gt_evr")
 #itemlist_end()
 #paragraph_end()
 #paragraph()
 Appstream module/stream filtering:
 #itemlist()
   #item("by module name, stream - field: module_stream; matcher: equals; value: modulaneme:stream")
 #itemlist_end()
 Note: Only 'allow' rule is supported for appstream filters.
 #paragraph_end()
 #paragraph()
 Note: The 'matches' matcher works on Java regular expressions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createFilter)
		},
	}

	cmd.Flags().String("Name", "", "Filter name")
	cmd.Flags().String("Rule", "", "Filter rule ('deny' or 'allow')")
	cmd.Flags().String("EntityType", "", "Filter entityType ('package' or 'erratum')")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func createFilter(globalFlags *types.GlobalFlags, flags *createFilterFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.Name, flags.Rule, flags.EntityType, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

