package master

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return the current default-Master for this Slave
func GetDefaultMaster(cnxDetails *api.ConnectionDetails) (*types.IssMaster, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "sync/master"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

	res, err := api.Get[types.IssMaster](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDefaultMaster: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
