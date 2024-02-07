package scap

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return a full list of RuleResults for given OpenSCAP XCCDF scan.
func GetXccdfScanRuleResults(cnxDetails *api.ConnectionDetails, Xid int) (*types.#return_array_begin()
   $XccdfRuleResultDtoSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system/scap"
	params := ""
	if Xid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
   $XccdfRuleResultDtoSerializer
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getXccdfScanRuleResults: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
