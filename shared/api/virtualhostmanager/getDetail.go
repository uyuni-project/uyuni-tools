package virtualhostmanager

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Gets details of a Virtual Host Manager with a given label
func GetDetail(cnxDetails *api.ConnectionDetails, Label string) (*types.VirtualHostManager, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "virtualhostmanager"
	params := ""
	if Label {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.VirtualHostManager](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetail: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
