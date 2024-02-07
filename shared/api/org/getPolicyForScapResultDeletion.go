package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the status of SCAP result deletion settings for the given
 organization.
func GetPolicyForScapResultDeletion(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#struct_begin("scap_deletion_info")
         #prop_desc("boolean", "enabled", "Deletion of SCAP results is enabled")
         #prop_desc("int", "retention_period",
             "Period (in days) after which a scan can be deleted (if enabled).")
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

    res, err := api.Get[types.#struct_begin("scap_deletion_info")
         #prop_desc("boolean", "enabled", "Deletion of SCAP results is enabled")
         #prop_desc("int", "retention_period",
             "Period (in days) after which a scan can be deleted (if enabled).")
     #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getPolicyForScapResultDeletion: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
