package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// get details about a ping sent to a system using OSA
func GetOsaPing(cnxDetails *api.ConnectionDetails, LoggedInUser User, Sid int) (*types.#struct_begin("osaPing")
          #prop_desc("string" "state"
          "state of the system (unknown, online, offline)")
          #prop_desc("$date" "lastMessageTime"
          "time of the last received response
          (1970/01/01 00:00:00 if never received a response)")
          #prop_desc("$date" "lastPingTime"
          "time of the last sent ping
          (1970/01/01 00:00:00 if no ping is pending")
      #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if LoggedInUser {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("osaPing")
          #prop_desc("string" "state"
          "state of the system (unknown, online, offline)")
          #prop_desc("$date" "lastMessageTime"
          "time of the last received response
          (1970/01/01 00:00:00 if never received a response)")
          #prop_desc("$date" "lastPingTime"
          "time of the last sent ping
          (1970/01/01 00:00:00 if no ping is pending")
      #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getOsaPing: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
