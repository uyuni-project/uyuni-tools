package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule package installation for several systems.
func SchedulePackageInstallByNevra(cnxDetails *api.ConnectionDetails, Sids []int, $param.getFlagName() $param.getType(), EarliestOccurrence $date, $param.getFlagName() $param.getType(), AllowModules bool, Sid int, AllowModules bool) (*types.#array_single("int", "actionId"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"$param.getName()":       $param.getFlagName(),
		"earliestOccurrence":       EarliestOccurrence,
		"$param.getName()":       $param.getFlagName(),
		"allowModules":       AllowModules,
		"sid":       Sid,
		"allowModules":       AllowModules,
	}

	res, err := api.Post[types.#array_single("int", "actionId")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schedulePackageInstallByNevra: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
