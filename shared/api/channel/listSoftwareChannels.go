package channel

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all visible software channels.
func ListSoftwareChannels(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
      #struct_begin("channel")
          #prop("string", "label")
          #prop("string", "name")
          #prop("string", "parent_label")
          #prop("string", "end_of_life")
          #prop("string", "arch")
      #struct_end()
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      #struct_begin("channel")
          #prop("string", "label")
          #prop("string", "name")
          #prop("string", "parent_label")
          #prop("string", "end_of_life")
          #prop("string", "arch")
      #struct_end()
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSoftwareChannels: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
