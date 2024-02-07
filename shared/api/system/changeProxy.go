package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Connect given systems to another proxy.
func ChangeProxy(cnxDetails *api.ConnectionDetails, Sids []int, ProxyId int) (*types.#array_single("int", "actionIds", "list of scheduled action ids"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"proxyId":       ProxyId,
	}

	res, err := api.Post[types.#array_single("int", "actionIds", "list of scheduled action ids")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute changeProxy: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
