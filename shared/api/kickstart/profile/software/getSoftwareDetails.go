package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Gets kickstart profile software details.
func GetSoftwareDetails(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#struct_begin("kickstart packages info")
              #prop_desc("string", "noBase", "install @Base package group")
              #prop_desc("string", "ignoreMissing", "ignore missing packages")
          #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile/software"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("kickstart packages info")
              #prop_desc("string", "noBase", "install @Base package group")
              #prop_desc("string", "ignoreMissing", "ignore missing packages")
          #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getSoftwareDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
