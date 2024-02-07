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
 should. This method lists the errata that will be updated if you run the
 syncErrata method.
func ListErrataNeedingSync(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#return_array_begin()
          $ErrataOverviewSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel/software"
	params := ""
	if ChannelLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          $ErrataOverviewSerializer
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listErrataNeedingSync: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
