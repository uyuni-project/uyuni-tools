package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of subscribable base channels.
func ListSubscribableBaseChannels(cnxDetails *api.ConnectionDetails, Sid int) (*types.#return_array_begin()
      #struct_begin("channel")
          #prop_desc("int" "id" "Base Channel ID.")
          #prop_desc("string" "name" "Name of channel.")
          #prop_desc("string" "label" "Label of Channel")
          #prop_desc("int", "current_base", "1 indicates it is the current base
                                      channel")
      #struct_end()
  #array_end(), error) {
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

    res, err := api.Get[types.#return_array_begin()
      #struct_begin("channel")
          #prop_desc("int" "id" "Base Channel ID.")
          #prop_desc("string" "name" "Name of channel.")
          #prop_desc("string" "label" "Label of Channel")
          #prop_desc("int", "current_base", "1 indicates it is the current base
                                      channel")
      #struct_end()
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSubscribableBaseChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
