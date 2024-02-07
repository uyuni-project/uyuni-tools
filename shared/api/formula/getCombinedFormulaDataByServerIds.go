package formula

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return the list of formulas a server and all his groups have.
func GetCombinedFormulaDataByServerIds(cnxDetails *api.ConnectionDetails, FormulaName string, Sids []int) (*types.#return_array_begin()
     $FormulaDataSerializer
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "formula"
	params := ""
	if FormulaName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if Sids {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     $FormulaDataSerializer
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getCombinedFormulaDataByServerIds: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
