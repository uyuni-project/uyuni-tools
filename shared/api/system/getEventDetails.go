package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the details of the event associated with the specified server and event.
             The event id must be a value returned by the system.getEventHistory API.
func GetEventDetails(cnxDetails *api.ConnectionDetails, Sid int, Eid int) (*types.#return_array_begin()
           $SystemEventDetailsDtoSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if Eid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
           $SystemEventDetailsDtoSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getEventDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
