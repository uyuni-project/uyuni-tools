package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Assign schedule with given name to systems with given IDs.
 Throws a PermissionCheckFailureException when some of the systems are not accessible by the user.
 Throws a InvalidParameterException when some of the systems have pending actions that are not allowed in the
 maintenance mode.
func AssignScheduleToSystems(cnxDetails *api.ConnectionDetails, ScheduleName string, $param.getFlagName() $param.getType(), $param.getFlagName() $param.getType()) (*types.#array_single("int", "number of involved systems"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"scheduleName":       ScheduleName,
		"$param.getName()":       $param.getFlagName(),
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#array_single("int", "number of involved systems")](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute assignScheduleToSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
