package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the date the system was registered.
func GetRegistrationDate(cnxDetails *api.ConnectionDetails, Sid int) (*types.#param_desc("$date", "date", "The date the system was registered,
 in local time"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("$date", "date", "The date the system was registered,
 in local time")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getRegistrationDate: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
