package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a Content Filter
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
 Note: The 'matches' matcher works on Java regular expressions.
func CreateFilter(cnxDetails *api.ConnectionDetails, Name string, Rule string, EntityType string, $param.getFlagName() $param.getType()) (*types.ContentFilter, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
		"rule":       Rule,
		"entityType":       EntityType,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.ContentFilter](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createFilter: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
