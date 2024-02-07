package admin

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new ssh connection data to extract data from
func Create(cnxDetails *api.ConnectionDetails, Description string, Host string, Port int, Username string, Password string, Key string, KeyPassword string, BastionHost string, BastionPort int, BastionUsername string, BastionPassword string, BastionKey string, BastionKeyPassword string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"description":       Description,
		"host":       Host,
		"port":       Port,
		"username":       Username,
		"password":       Password,
		"key":       Key,
		"keyPassword":       KeyPassword,
		"bastionHost":       BastionHost,
		"bastionPort":       BastionPort,
		"bastionUsername":       BastionUsername,
		"bastionPassword":       BastionPassword,
		"bastionKey":       BastionKey,
		"bastionKeyPassword":       BastionKeyPassword,
	}

	res, err := api.Post[types.#return_int_success()](client, "admin", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
