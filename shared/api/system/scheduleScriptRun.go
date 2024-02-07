package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule a script to run.
func ScheduleScriptRun(cnxDetails *api.ConnectionDetails, Label string, $param.getFlagName() $param.getType(), Username string, Groupname string, Timeout int, Script string, EarliestOccurrence $date, Sid int) (*types.#param_desc("int", "id", "ID of the script run action created. Can be used to fetch
 results with system.getScriptResults"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"$param.getName()":       $param.getFlagName(),
		"username":       Username,
		"groupname":       Groupname,
		"timeout":       Timeout,
		"script":       Script,
		"earliestOccurrence":       EarliestOccurrence,
		"sid":       Sid,
	}

	res, err := api.Post[types.#param_desc("int", "id", "ID of the script run action created. Can be used to fetch
 results with system.getScriptResults")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleScriptRun: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
