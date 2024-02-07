package channel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all retired software channels.  These are channels that the user's
 organization is entitled to, but are no longer supported because they have reached
 their 'end-of-life' date.
func ListRetiredChannels(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
         $ChannelTreeNodeSerializer
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
         $ChannelTreeNodeSerializer
     #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listRetiredChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
