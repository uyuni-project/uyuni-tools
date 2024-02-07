package formula

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Return the list of formulas a server group has.
func GetFormulasByGroupId(cnxDetails *api.ConnectionDetails, SystemGroupId int) (*types.#array_single("string", "the list of formulas"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "formula"
	params := ""
	if SystemGroupId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "the list of formulas")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getFormulasByGroupId: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
