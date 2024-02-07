package virtualhostmanager

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List all available modules from virtual-host-gatherer
func ListAvailableVirtualHostGathererModules(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "moduleName"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "virtualhostmanager"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "moduleName")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listAvailableVirtualHostGathererModules: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
