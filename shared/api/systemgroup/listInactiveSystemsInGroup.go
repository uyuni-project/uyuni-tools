package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists inactive systems within a server group using a
          specified inactivity time.
func ListInactiveSystemsInGroup(cnxDetails *api.ConnectionDetails, SystemGroupName string, DaysInactive int) (*types.#array_single("int", "server_id"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "systemgroup"
	params := ""
	if SystemGroupName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if DaysInactive {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("int", "server_id")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listInactiveSystemsInGroup: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
