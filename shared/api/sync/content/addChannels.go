package content

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add a new channel to the #product() database
func AddChannels(cnxDetails *api.ConnectionDetails, ChannelLabel string, MirrorUrl string) (*types.#array_single("string", "enabled channel labels"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"mirrorUrl":       MirrorUrl,
	}

	res, err := api.Post[types.#array_single("string", "enabled channel labels")](client, "sync/content", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
