package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Retract schedule with given name from systems with given IDs
 Throws a PermissionCheckFailureException when some of the systems are not accessible by the user.
func RetractScheduleFromSystems(cnxDetails *api.ConnectionDetails, $param.getFlagName() $param.getType()) (*types.#array_single("int", "number of involved systems"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#array_single("int", "number of involved systems")](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute retractScheduleFromSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
