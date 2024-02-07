package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get system name and last check in information for the given system ID.
func GetName(cnxDetails *api.ConnectionDetails, Sid string) (*types.#struct_begin("name info")
      #prop_desc("int", "id", "Server id")
      #prop_desc("string", "name", "Server name")
      #prop_desc("$date", "last_checkin", "Last time server
              successfully checked in")
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

    res, err := api.Get[types.#struct_begin("name info")
      #prop_desc("int", "id", "Server id")
      #prop_desc("string", "name", "Server name")
      #prop_desc("$date", "last_checkin", "Last time server
              successfully checked in")
  #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getName: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
