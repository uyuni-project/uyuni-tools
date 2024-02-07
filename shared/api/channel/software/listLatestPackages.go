package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists the packages with the latest version (including release and
 epoch) for the given channel
func ListLatestPackages(cnxDetails *api.ConnectionDetails, ChannelLabel string) (*types.#return_array_begin()
          #struct_begin("package")
              #prop("string", "name")
              #prop("string", "version")
              #prop("string", "release")
              #prop("string", "epoch")
              #prop("int", "id")
              #prop("string", "arch_label")
          #struct_end()
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "channel/software"
	params := ""
	if ChannelLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          #struct_begin("package")
              #prop("string", "name")
              #prop("string", "version")
              #prop("string", "release")
              #prop("string", "epoch")
              #prop("int", "id")
              #prop("string", "arch_label")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listLatestPackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
