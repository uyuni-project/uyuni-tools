package formula

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set the formulas of a server group.
func SetFormulasOfGroup(cnxDetails *api.ConnectionDetails, SystemGroupId int, Formulas []string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"systemGroupId":       SystemGroupId,
		"formulas":       Formulas,
	}

	res, err := api.Post[types.#return_int_success()](client, "formula", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setFormulasOfGroup: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
