package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Import an image and schedule an inspect afterwards
func ImportContainerImage(cnxDetails *api.ConnectionDetails, Name string, Version string, BuildHostId int, StoreLabel string, ActivationKey string, EarliestOccurrence $date) (*types.#param("int", "the ID of the inspect action created"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
		"version":       Version,
		"buildHostId":       BuildHostId,
		"storeLabel":       StoreLabel,
		"activationKey":       ActivationKey,
		"earliestOccurrence":       EarliestOccurrence,
	}

	res, err := api.Post[types.#param("int", "the ID of the inspect action created")](client, "image", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute importContainerImage: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
