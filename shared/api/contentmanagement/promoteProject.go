package contentmanagement

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Promote an Environment in a Project
func PromoteProject(cnxDetails *api.ConnectionDetails, ProjectLabel string, EnvLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"projectLabel":       ProjectLabel,
		"envLabel":       EnvLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "contentmanagement", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute promoteProject: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
