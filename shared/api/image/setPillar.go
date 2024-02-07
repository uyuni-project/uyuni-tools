package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set pillar data of an image. The "size" entries should be passed as string.
func SetPillar(cnxDetails *api.ConnectionDetails, ImageId int, PillarData struct) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"imageId":       ImageId,
		"pillarData":       PillarData,
	}

	res, err := api.Post[types.#return_int_success()](client, "image", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setPillar: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
