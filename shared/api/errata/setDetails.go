package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set erratum details. All arguments are optional and will only be modified
 if included in the struct. This method will only allow for modification of custom
 errata created either through the UI or API.
func SetDetails(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"advisoryName":       AdvisoryName,
	}

	res, err := api.Post[types.#return_int_success()](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
