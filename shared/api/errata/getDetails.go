package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Retrieves the details for the erratum matching the given advisory name.
func GetDetails(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#struct_begin("erratum")
          #prop("int", "id")
          #prop("string", "issue_date")
          #prop("string", "update_date")
          #prop_desc("string", "last_modified_date", "last time the erratum was modified.")
          #prop("string", "synopsis")
          #prop("int", "release")
          #prop("string", "advisory_status")
          #prop("string", "vendor_advisory")
          #prop("string", "type")
          #prop("string", "product")
          #prop("string", "errataFrom")
          #prop("string", "topic")
          #prop("string", "description")
          #prop("string", "references")
          #prop("string", "notes")
          #prop("string", "solution")
          #prop_desc("boolean", "reboot_suggested", "A boolean flag signaling whether a system reboot is
          advisable following the application of the errata. Typical example is upon kernel update.")
          #prop_desc("boolean", "restart_suggested", "A boolean flag signaling a weather reboot of
          the package manager is advisable following the application of the errata. This is commonly
          used to address update stack issues before proceeding with other updates.")
     #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "errata"
	params := ""
	if AdvisoryName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("erratum")
          #prop("int", "id")
          #prop("string", "issue_date")
          #prop("string", "update_date")
          #prop_desc("string", "last_modified_date", "last time the erratum was modified.")
          #prop("string", "synopsis")
          #prop("int", "release")
          #prop("string", "advisory_status")
          #prop("string", "vendor_advisory")
          #prop("string", "type")
          #prop("string", "product")
          #prop("string", "errataFrom")
          #prop("string", "topic")
          #prop("string", "description")
          #prop("string", "references")
          #prop("string", "notes")
          #prop("string", "solution")
          #prop_desc("boolean", "reboot_suggested", "A boolean flag signaling whether a system reboot is
          advisable following the application of the errata. Typical example is upon kernel update.")
          #prop_desc("boolean", "restart_suggested", "A boolean flag signaling a weather reboot of
          the package manager is advisable following the application of the errata. This is commonly
          used to address update stack issues before proceeding with other updates.")
     #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
