package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule a comparison of the latest revision of a file
 against the version deployed on a list of systems.
func ScheduleFileComparisons(cnxDetails *api.ConnectionDetails, Label string, Path string, Sids []long) (*types.#param_desc("int", "actionId", "the action ID of the scheduled action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"path":       Path,
		"sids":       Sids,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "the action ID of the scheduled action")](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleFileComparisons: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
