package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Import an image and schedule an inspect afterwards
func ImportOSImage(cnxDetails *api.ConnectionDetails, Name string, Version string, Arch string) (*types.#param("int", "the ID of the image"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
		"version":       Version,
		"arch":       Arch,
	}

	res, err := api.Post[types.#param("int", "the ID of the image")](client, "image", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute importOSImage: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
