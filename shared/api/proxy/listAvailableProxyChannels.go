package proxy

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List available version of proxy channel for system
 identified by the given client certificate i.e. systemid file.
func ListAvailableProxyChannels(cnxDetails *api.ConnectionDetails, Clientcert string) (*types.#array_single ("string", "version"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "proxy"
	params := ""
	if Clientcert {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single ("string", "version")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAvailableProxyChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
