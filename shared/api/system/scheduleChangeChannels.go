package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule an action to change the channels of the given system. Works for both traditional
 and Salt systems.
 This method accepts labels for the base and child channels.
 If the user provides an empty string for the channelLabel, the current base channel and
 all child channels will be removed from the system.
func ScheduleChangeChannels(cnxDetails *api.ConnectionDetails, Sid int, BaseChannelLabel string, ChildLabels []string, EarliestOccurrence $type, Sids []int) (*types.#param_desc("int", "id", "ID of the action scheduled, otherwise exception thrown
 on error"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"baseChannelLabel":       BaseChannelLabel,
		"childLabels":       ChildLabels,
		"earliestOccurrence":       EarliestOccurrence,
		"sids":       Sids,
	}

	res, err := api.Post[types.#param_desc("int", "id", "ID of the action scheduled, otherwise exception thrown
 on error")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleChangeChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
