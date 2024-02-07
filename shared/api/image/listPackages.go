package image

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the installed packages on the given image
func ListPackages(cnxDetails *api.ConnectionDetails, ImageId int) (*types.#return_array_begin()
          #struct_begin("package")
                 #prop("string", "name")
                 #prop("string", "version")
                 #prop("string", "release")
                 #prop("string", "epoch")
                 #prop("string", "arch")
          #struct_end()
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "image"
	params := ""
	if ImageId {
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
                 #prop("string", "arch")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listPackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
