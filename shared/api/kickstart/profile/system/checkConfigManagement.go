package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Check the configuration management status for a kickstart profile.
func CheckConfigManagement(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#param("boolean", "true if configuration management is enabled; otherwise, false"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
	}

	res, err := api.Post[types.#param("boolean", "true if configuration management is enabled; otherwise, false")](client, "kickstart/profile/system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute checkConfigManagement: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
