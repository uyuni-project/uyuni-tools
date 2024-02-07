package channel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the most popular software channels.  Channels that have at least
 the number of systems subscribed as specified by the popularity count will be
 returned.
func ListPopularChannels(cnxDetails *api.ConnectionDetails, PopularityCount int) (*types.#return_array_begin()
         $ChannelTreeNodeSerializer
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel"
	params := ""
	if PopularityCount {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
         $ChannelTreeNodeSerializer
     #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listPopularChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
