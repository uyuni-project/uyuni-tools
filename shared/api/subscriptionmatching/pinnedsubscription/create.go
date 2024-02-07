package pinnedsubscription

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Creates a Pinned Subscription based on given subscription and system
func Create(cnxDetails *api.ConnectionDetails, SubscriptionId int, Sid int) (*types.PinnedSubscription, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"subscriptionId":       SubscriptionId,
		"sid":       Sid,
	}

	res, err := api.Post[types.PinnedSubscription](client, "subscriptionmatching/pinnedsubscription", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
