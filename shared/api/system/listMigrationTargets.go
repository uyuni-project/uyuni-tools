package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List possible migration targets for a system
func ListMigrationTargets(cnxDetails *api.ConnectionDetails, Sid int, ExcludeTargetWhereMissingSuccessors bool) (*types.#return_array_begin()
          #struct_begin("migrationtarget")
                 #prop("string", "ident")
                 #prop("string", "friendly")
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
	if ExcludeTargetWhereMissingSuccessors {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          #struct_begin("migrationtarget")
                 #prop("string", "ident")
                 #prop("string", "friendly")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listMigrationTargets: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
