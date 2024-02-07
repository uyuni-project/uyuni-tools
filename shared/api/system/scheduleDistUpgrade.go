package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule a dist upgrade for a system. This call takes a list of channel
 labels that the system will be subscribed to before performing the dist upgrade.
 Note: You can seriously damage your system with this call, use it only if you really
 know what you are doing! Make sure that the list of channel labels is complete and in
 any case do a dry run before scheduling an actual dist upgrade.
func ScheduleDistUpgrade(cnxDetails *api.ConnectionDetails, Sid int, Channels []string, DryRun bool, EarliestOccurrence $date, AllowVendorChange bool) (*types.#param("int", "actionId", "The action id of the scheduled action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"channels":       Channels,
		"dryRun":       DryRun,
		"earliestOccurrence":       EarliestOccurrence,
		"allowVendorChange":       AllowVendorChange,
	}

	res, err := api.Post[types.#param("int", "actionId", "The action id of the scheduled action")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleDistUpgrade: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
