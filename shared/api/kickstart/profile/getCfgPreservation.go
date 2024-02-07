package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get ks.cfg preservation option for a kickstart profile.
func GetCfgPreservation(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#param_desc("boolean", "preserve", "The value of the option.
      True means that ks.cfg will be copied to /root, false means that it will not"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("boolean", "preserve", "The value of the option.
      True means that ks.cfg will be copied to /root, false means that it will not")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getCfgPreservation: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
