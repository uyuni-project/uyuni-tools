package actionchain

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List currently available action chains.
func ListChains(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
                      #struct_begin("chain")
                        #prop_desc("string", "label", "Label of an Action Chain")
                        #prop_desc("string", "entrycount",
                                   "Number of entries in the Action Chain")
                      #struct_end()
                    #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "actionchain"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
                      #struct_begin("chain")
                        #prop_desc("string", "label", "Label of an Action Chain")
                        #prop_desc("string", "entrycount",
                                   "Number of entries in the Action Chain")
                      #struct_end()
                    #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listChains: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
