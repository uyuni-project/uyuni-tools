package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule an image build
func ScheduleImageBuild(cnxDetails *api.ConnectionDetails, ProfileLabel string, Version string, BuildHostId int, EarliestOccurrence $date) (*types.#param("int", "the ID of the build action created"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"profileLabel":       ProfileLabel,
		"version":       Version,
		"buildHostId":       BuildHostId,
		"earliestOccurrence":       EarliestOccurrence,
	}

	res, err := api.Post[types.#param("int", "the ID of the build action created")](client, "image", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleImageBuild: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
