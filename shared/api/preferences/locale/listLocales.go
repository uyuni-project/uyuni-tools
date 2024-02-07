package locale

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of all understood locales. Can be
 used as input to setLocale.
func ListLocales(cnxDetails *api.ConnectionDetails) (*types.#array_single("string", "Locale code."), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "preferences/locale"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "Locale code.")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listLocales: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
