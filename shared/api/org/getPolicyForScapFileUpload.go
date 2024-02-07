package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the status of SCAP detailed result file upload settings
 for the given organization.
func GetPolicyForScapFileUpload(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#struct_begin("scap_upload_info")
         #prop_desc("boolean", "enabled",
             "Aggregation of detailed SCAP results is enabled.")
         #prop_desc("int", "size_limit",
             "Limit (in Bytes) for a single SCAP file upload.")
     #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "org"
	params := ""
	if OrgId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("scap_upload_info")
         #prop_desc("boolean", "enabled",
             "Aggregation of detailed SCAP results is enabled.")
         #prop_desc("int", "size_limit",
             "Limit (in Bytes) for a single SCAP file upload.")
     #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getPolicyForScapFileUpload: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
