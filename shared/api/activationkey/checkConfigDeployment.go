package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Check configuration file deployment status for the
 activation key specified.
func CheckConfigDeployment(cnxDetails *api.ConnectionDetails, Key string) (*types.#param_desc("int", "status", "1 if enabled, 0 if disabled, exception thrown otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"key":       Key,
	}

	res, err := api.Post[types.#param_desc("int", "status", "1 if enabled, 0 if disabled, exception thrown otherwise")](client, "activationkey", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute checkConfigDeployment: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
