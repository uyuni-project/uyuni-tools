package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Align the metadata of a channel to another channel.
func AlignMetadata(cnxDetails *api.ConnectionDetails, ChannelFromLabel string, ChannelToLabel string, MetadataType string) (*types.#param_desc("int", "result code", "1 when metadata has been aligned, 0 otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelFromLabel":       ChannelFromLabel,
		"channelToLabel":       ChannelToLabel,
		"metadataType":       MetadataType,
	}

	res, err := api.Post[types.#param_desc("int", "result code", "1 when metadata has been aligned, 0 otherwise")](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute alignMetadata: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
