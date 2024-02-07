package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the System Currency scores for all servers the user has access to
func GetSystemCurrencyScores(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
          #struct_begin("system currency")
              #prop("int", "sid")
              #prop("int", "critical security errata count")
              #prop("int", "important security errata count")
              #prop("int", "moderate security errata count")
              #prop("int", "low security errata count")
              #prop("int", "bug fix errata count")
              #prop("int", "enhancement errata count")
              #prop("int", "system currency score")
          #struct_end()
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
          #struct_begin("system currency")
              #prop("int", "sid")
              #prop("int", "critical security errata count")
              #prop("int", "important security errata count")
              #prop("int", "moderate security errata count")
              #prop("int", "low security errata count")
              #prop("int", "bug fix errata count")
              #prop("int", "enhancement errata count")
              #prop("int", "system currency score")
          #struct_end()
      #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getSystemCurrencyScores: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
