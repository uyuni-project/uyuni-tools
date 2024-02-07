package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule an action of a guest's host, to set that guest's CPU
          allocation
func SetGuestCpus(cnxDetails *api.ConnectionDetails, Sid int, NumOfCpus int) (*types.#param_desc("int", "actionID", "the action Id for the schedule action
              on the host system"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"numOfCpus":       NumOfCpus,
	}

	res, err := api.Post[types.#param_desc("int", "actionID", "the action Id for the schedule action
              on the host system")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setGuestCpus: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
