package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Clone an existing activation key.
func Clone(cnxDetails *api.ConnectionDetails, Key string, CloneDescription string) (*types.#param("string", "The new activation key"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"key":       Key,
		"cloneDescription":       CloneDescription,
	}

	res, err := api.Post[types.#param("string", "The new activation key")](client, "activationkey", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute clone: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
