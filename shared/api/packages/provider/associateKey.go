package provider

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Associate a package security key and with the package provider.
      If the provider or key doesn't exist, it is created. User executing the
      request must be a #product() administrator.
func AssociateKey(cnxDetails *api.ConnectionDetails, ProviderName string, Key string, Type string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"providerName":       ProviderName,
		"key":       Key,
		"type":       Type,
	}

	res, err := api.Post[types.#return_int_success()](client, "packages/provider", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute associateKey: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
