package formula

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set the formula form for the specified server.
func SetSystemFormulaData(cnxDetails *api.ConnectionDetails, SystemId int, FormulaName string, Content struct) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemId":       SystemId,
		"formulaName":       FormulaName,
		"content":       Content,
	}

	res, err := api.Post[types.#return_int_success()](client, "formula", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setSystemFormulaData: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
