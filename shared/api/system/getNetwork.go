package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the addresses and hostname for a given server.
func GetNetwork(cnxDetails *api.ConnectionDetails, Sid int) (*types.#struct_begin("network info")
              #prop_desc("string", "ip", "IPv4 address of server")
              #prop_desc("string", "ip6", "IPv6 address of server")
              #prop_desc("string", "hostname", "Hostname of server")
          #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("network info")
              #prop_desc("string", "ip", "IPv4 address of server")
              #prop_desc("string", "ip6", "IPv6 address of server")
              #prop_desc("string", "hostname", "Hostname of server")
          #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getNetwork: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
