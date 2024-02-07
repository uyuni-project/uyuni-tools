package configchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return a list of files in a channel.
func ListFiles(cnxDetails *api.ConnectionDetails, Label string) (*types.#return_array_begin()
 $ConfigFileDtoSerializer
 #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "configchannel"
	params := ""
	if Label {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
 $ConfigFileDtoSerializer
 #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listFiles: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
