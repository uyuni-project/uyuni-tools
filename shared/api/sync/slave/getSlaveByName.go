package slave

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Find a Slave by specifying its Fully-Qualified Domain Name
func GetSlaveByName(cnxDetails *api.ConnectionDetails, SlaveFqdn string) (*types.IssSlave, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "sync/slave"
	params := ""
	if SlaveFqdn {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.IssSlave](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getSlaveByName: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
