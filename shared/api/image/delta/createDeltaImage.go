package delta

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Import an image and schedule an inspect afterwards. The "size" entries in the pillar
 should be passed as string.
func CreateDeltaImage(cnxDetails *api.ConnectionDetails, SourceImageId int, TargetImageId int, File string, Pillar struct) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sourceImageId":       SourceImageId,
		"targetImageId":       TargetImageId,
		"file":       File,
		"pillar":       Pillar,
	}

	res, err := api.Post[types.#return_int_success()](client, "image/delta", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createDeltaImage: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
