package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove a source package.
func RemoveSourcePackage(cnxDetails *api.ConnectionDetails, Psid int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"psid":       Psid,
	}

	res, err := api.Post[types.#return_int_success()](client, "packages", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeSourcePackage: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
