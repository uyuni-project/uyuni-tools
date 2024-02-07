package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Obtains a reactivation key for this server.
func ObtainReactivationKey(cnxDetails *api.ConnectionDetails, Sid int, ClientCert string) (*types.#param("string", "key"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"clientCert":       ClientCert,
	}

	res, err := api.Post[types.#param("string", "key")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute obtainReactivationKey: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
