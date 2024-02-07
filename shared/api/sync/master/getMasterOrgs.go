package master

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all organizations the specified Master has exported to this Slave
func GetMasterOrgs(cnxDetails *api.ConnectionDetails, MasterId int) (*types.#return_array_begin()
     $IssMasterOrgSerializer
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "sync/master"
	params := ""
	if MasterId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     $IssMasterOrgSerializer
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getMasterOrgs: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
