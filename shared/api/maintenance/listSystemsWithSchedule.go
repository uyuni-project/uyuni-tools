package maintenance

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List IDs of systems that have given schedule assigned
 Throws a PermissionCheckFailureException when some of the systems are not accessible by the user.
func ListSystemsWithSchedule(cnxDetails *api.ConnectionDetails, ScheduleName string) (*types.#array_single("int", "system IDs"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"scheduleName":       ScheduleName,
	}

	res, err := api.Post[types.#array_single("int", "system IDs")](client, "maintenance", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSystemsWithSchedule: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
