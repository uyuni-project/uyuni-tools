package distchannel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Sets, overrides (/removes if channelLabel empty)
 a distribution channel map within an organization
func SetMapForOrg(cnxDetails *api.ConnectionDetails, Os string, Release string, ArchName string, ChannelLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"os":       Os,
		"release":       Release,
		"archName":       ArchName,
		"channelLabel":       ChannelLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "distchannel", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setMapForOrg: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
