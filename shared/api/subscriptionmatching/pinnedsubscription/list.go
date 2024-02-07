package pinnedsubscription

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists all PinnedSubscriptions
func List(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
         $PinnedSubscriptionSerializer
     #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#return_array_begin()
         $PinnedSubscriptionSerializer
     #array_end()](client, "subscriptionmatching/pinnedsubscription", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute list: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
