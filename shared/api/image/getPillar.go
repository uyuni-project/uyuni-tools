package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get pillar data of an image. The "size" entries are converted to string.
func GetPillar(cnxDetails *api.ConnectionDetails, ImageId int) (*types.#param("struct", "the pillar data"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "image"
	params := ""
	if ImageId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#param("struct", "the pillar data")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getPillar: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
