package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Synchronize all files on the disk to the current state of the database.
func SyncSaltFilesOnDisk(cnxDetails *api.ConnectionDetails, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_int_success()](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute syncSaltFilesOnDisk: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
