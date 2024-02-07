package audit

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List visible images with their patch status regarding a given CVE
 identifier. Please note that the query code relies on data that is pre-generated
 by the 'cve-server-channels' taskomatic job.
func ListImagesByPatchStatus(cnxDetails *api.ConnectionDetails, CveIdentifier string) (*types.#return_array_begin() $CVEAuditImageSerializer #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "audit"
	params := ""
	if CveIdentifier {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin() $CVEAuditImageSerializer #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listImagesByPatchStatus: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
