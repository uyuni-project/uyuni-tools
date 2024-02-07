package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a custom errata
func Create(cnxDetails *api.ConnectionDetails, $param.getFlagName() $param.getType(), $param.getFlagName() $param.getType(), PackageIds []int, $param.getFlagName() $param.getType()) (*types.Errata, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"$param.getName()":       $param.getFlagName(),
		"$param.getName()":       $param.getFlagName(),
		"packageIds":       PackageIds,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.Errata](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
