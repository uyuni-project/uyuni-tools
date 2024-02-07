package keys

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// returns all the data associated with the given key
func GetDetails(cnxDetails *api.ConnectionDetails, Description string) (*types.#struct_begin("key")
          #prop("string", "description")
          #prop("string", "type")
          #prop("string", "content")
      #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/keys"
	params := ""
	if Description {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("key")
          #prop("string", "description")
          #prop("string", "type")
          #prop("string", "content")
      #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
