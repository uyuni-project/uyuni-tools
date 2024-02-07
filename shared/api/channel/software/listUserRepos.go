package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of ContentSource (repos) that the user can see
func ListUserRepos(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
          #struct_begin("map")
              #prop_desc("long","id", "ID of the repo")
              #prop_desc("string","label", "label of the repo")
              #prop_desc("string","sourceUrl", "URL of the repo")
          #struct_end()
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel/software"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          #struct_begin("map")
              #prop_desc("long","id", "ID of the repo")
              #prop_desc("string","label", "label of the repo")
              #prop_desc("string","sourceUrl", "URL of the repo")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listUserRepos: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
