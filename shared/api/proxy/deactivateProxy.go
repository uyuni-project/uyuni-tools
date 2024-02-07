package proxy

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Deactivates the proxy identified by the given client
 certificate i.e. systemid file.
func DeactivateProxy(cnxDetails *api.ConnectionDetails, Clientcert string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"clientcert":       Clientcert,
	}

	res, err := api.Post[types.#return_int_success()](client, "proxy", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deactivateProxy: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
