package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set contact/support information for given channel.
func SetContactDetails(cnxDetails *api.ConnectionDetails, ChannelLabel string, MaintainerName string, MaintainerEmail string, MaintainerPhone string, SupportPolicy string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"maintainerName":       MaintainerName,
		"maintainerEmail":       MaintainerEmail,
		"maintainerPhone":       MaintainerPhone,
		"supportPolicy":       SupportPolicy,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setContactDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
