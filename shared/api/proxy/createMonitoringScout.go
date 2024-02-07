package proxy

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create Monitoring Scout for proxy.
func CreateMonitoringScout(cnxDetails *api.ConnectionDetails, Clientcert string) (*types.#param("string", ""), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"clientcert":       Clientcert,
	}

	res, err := api.Post[types.#param("string", "")](client, "proxy", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createMonitoringScout: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
