package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Trigger immediate repo synchronization
func SyncRepo(cnxDetails *api.ConnectionDetails, ChannelLabels []string, ChannelLabel string, $param.getFlagName() $param.getType(), CronExpr string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabels":       ChannelLabels,
		"channelLabel":       ChannelLabel,
		"$param.getName()":       $param.getFlagName(),
		"cronExpr":       CronExpr,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute syncRepo: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
