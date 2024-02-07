package systemgroup

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Subscribe given config channels to a system group
func SubscribeConfigChannel(cnxDetails *api.ConnectionDetails, SystemGroupName string, ConfigChannelLabels []string) (*types.1 on success, exception on failure, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemGroupName":       SystemGroupName,
		"configChannelLabels":       ConfigChannelLabels,
	}

	res, err := api.Post[types.1 on success, exception on failure](client, "systemgroup", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute subscribeConfigChannel: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
