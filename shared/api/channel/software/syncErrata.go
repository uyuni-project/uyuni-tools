package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// If you have synced a new channel then patches
 will have been updated with the packages that are in the newly
 synced channel. A cloned erratum will not have been automatically updated
 however. If you cloned a channel that includes those cloned errata and
 should include the new packages, they will not be included when they
 should. This method updates all the errata in the given cloned channel
 with packages that have recently been added, and ensures that all the
 packages you expect are in the channel. It also updates cloned errata
 attributes like advisoryStatus.
func SyncErrata(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute syncErrata: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
