package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule the Action Chain so that its actions will actually occur.
func ScheduleChain(cnxDetails *api.ConnectionDetails, ChainLabel string, Date $date) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"chainLabel":       ChainLabel,
		"date":       Date,
	}

	res, err := api.Post[types.#return_int_success()](client, "actionchain", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleChain: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
