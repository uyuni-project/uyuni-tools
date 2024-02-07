package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete image file
func AddImageFile(cnxDetails *api.ConnectionDetails, ImageId int, File string, Type string, External bool) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"imageId":       ImageId,
		"file":       File,
		"type":       Type,
		"external":       External,
	}

	res, err := api.Post[types.#return_int_success()](client, "image", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addImageFile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
