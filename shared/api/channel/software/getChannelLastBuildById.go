package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the last build date of the repomd.xml file
 for the given channel as a localised string.
func GetChannelLastBuildById(cnxDetails *api.ConnectionDetails, Id int) (*types.#param_desc("date", "date", "the last build date of the repomd.xml file as a localised string"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel/software"
	params := ""
	if Id {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param_desc("date", "date", "the last build date of the repomd.xml file as a localised string")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getChannelLastBuildById: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
