package slave

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get all orgs this Master is willing to export to the specified Slave
func GetAllowedOrgs(cnxDetails *api.ConnectionDetails, SlaveId int) (*types.#array_single("int", "ids of allowed organizations"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "sync/slave"
	params := ""
	if SlaveId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("int", "ids of allowed organizations")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getAllowedOrgs: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
