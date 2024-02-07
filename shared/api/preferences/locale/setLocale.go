package locale

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set a user's locale.
func SetLocale(cnxDetails *api.ConnectionDetails, Login string, Locale string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"login":       Login,
		"locale":       Locale,
	}

	res, err := api.Post[types.#return_int_success()](client, "preferences/locale", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setLocale: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
