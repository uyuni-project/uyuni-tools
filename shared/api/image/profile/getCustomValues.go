package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the custom data values defined for the image profile
func GetCustomValues(cnxDetails *api.ConnectionDetails, Label string) (*types.#struct_begin("the map of custom labels to custom values")
      #prop("string", "custom info label")
      #prop("string", "value")
    #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "image/profile"
	params := ""
	if Label {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("the map of custom labels to custom values")
      #prop("string", "custom info label")
      #prop("string", "value")
    #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getCustomValues: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
