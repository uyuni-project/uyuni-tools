package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule an immediate configuration deployment for all systems
    subscribed to a particular configuration channel.
func DeployAllSystems(cnxDetails *api.ConnectionDetails, Label string, Date $date, FilePath string, Date $date) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"date":       Date,
		"filePath":       FilePath,
		"date":       Date,
	}

	res, err := api.Post[types.#return_int_success()](client, "configchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deployAllSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
