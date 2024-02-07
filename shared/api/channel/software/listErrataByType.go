package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the errata of a specific type that are applicable to a channel
func ListErrataByType(cnxDetails *api.ConnectionDetails, ChannelLabel string, AdvisoryType string) (*types.#return_array_begin()
          #struct_begin("errata")
              #prop_desc("string","advisory", "name of the advisory")
              #prop_desc("string","issue_date",
                         "date format follows YYYY-MM-DD HH24:MI:SS")
              #prop_desc("string","update_date",
                         "date format follows YYYY-MM-DD HH24:MI:SS")
              #prop("string","synopsis")
              #prop("string","advisory_type")
              #prop_desc("string","last_modified_date",
                         "date format follows YYYY-MM-DD HH24:MI:SS")
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
	if AdvisoryType {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          #struct_begin("errata")
              #prop_desc("string","advisory", "name of the advisory")
              #prop_desc("string","issue_date",
                         "date format follows YYYY-MM-DD HH24:MI:SS")
              #prop_desc("string","update_date",
                         "date format follows YYYY-MM-DD HH24:MI:SS")
              #prop("string","synopsis")
              #prop("string","advisory_type")
              #prop_desc("string","last_modified_date",
                         "date format follows YYYY-MM-DD HH24:MI:SS")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listErrataByType: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
