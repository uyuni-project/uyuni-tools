package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the addresses and hostname for a given list of systems.
func GetNetworkForSystems(cnxDetails *api.ConnectionDetails, Sids []int) (*types.#return_array_begin()
     #struct_begin("network info")
       #prop_desc("int", "system_id", "ID of the system")
       #prop_desc("string", "ip", "IPv4 address of system")
       #prop_desc("string", "ip6", "IPv6 address of system")
       #prop_desc("string", "hostname", "Hostname of system")
       #prop_desc("string", "primary_fqdn", "Primary FQDN of system")
     #struct_end()
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sids {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     #struct_begin("network info")
       #prop_desc("int", "system_id", "ID of the system")
       #prop_desc("string", "ip", "IPv4 address of system")
       #prop_desc("string", "ip6", "IPv6 address of system")
       #prop_desc("string", "hostname", "Hostname of system")
       #prop_desc("string", "primary_fqdn", "Primary FQDN of system")
     #struct_end()
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getNetworkForSystems: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
