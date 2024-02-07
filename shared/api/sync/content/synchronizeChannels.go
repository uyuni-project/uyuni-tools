package content

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// (Deprecated) Synchronize channels between the Customer Center
             and the #product() database.
func SynchronizeChannels(cnxDetails *api.ConnectionDetails, MirrorUrl string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"mirrorUrl":       MirrorUrl,
	}

	res, err := api.Post[types.#return_int_success()](client, "sync/content", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute synchronizeChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
