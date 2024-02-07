package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule package installation for several systems.
func SchedulePackageInstall(cnxDetails *api.ConnectionDetails, Sids []int, PackageIds []int, EarliestOccurrence $date, AllowModules bool, Sid int) (*types.#array_single("int", "actionId"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"packageIds":       PackageIds,
		"earliestOccurrence":       EarliestOccurrence,
		"allowModules":       AllowModules,
		"sid":       Sid,
	}

	res, err := api.Post[types.#array_single("int", "actionId")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schedulePackageInstall: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
