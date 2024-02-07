package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all actions in the particular Action Chain.
func ListChainActions(cnxDetails *api.ConnectionDetails, ChainLabel string) (*types.#return_array_begin()
                      #struct_begin("entry")
                        #prop_desc("int", "id", "Action ID")
                        #prop_desc("string", "label", "Label of an Action")
                        #prop_desc("string", "created", "Created date/time")
                        #prop_desc("string", "earliest",
                                   "Earliest scheduled date/time")
                        #prop_desc("string", "type", "Type of the action")
                        #prop_desc("string", "modified", "Modified date/time")
                        #prop_desc("string", "cuid", "Creator UID")
                      #struct_end()
                    #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "actionchain"
	params := ""
	if ChainLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
                      #struct_begin("entry")
                        #prop_desc("int", "id", "Action ID")
                        #prop_desc("string", "label", "Label of an Action")
                        #prop_desc("string", "created", "Created date/time")
                        #prop_desc("string", "earliest",
                                   "Earliest scheduled date/time")
                        #prop_desc("string", "type", "Type of the action")
                        #prop_desc("string", "modified", "Modified date/time")
                        #prop_desc("string", "cuid", "Creator UID")
                      #struct_end()
                    #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listChainActions: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
