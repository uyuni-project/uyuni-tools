package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Retrieves the locale for a kickstart profile.
func GetLocale(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#struct_begin("locale info")
              #prop("string", "locale")
              #prop("boolean", "useUtc")
                  #options()
                      #item_desc ("true", "the hardware clock uses UTC")
                      #item_desc ("false", "the hardware clock does not use UTC")
                  #options_end()
          #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile/system"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("locale info")
              #prop("string", "locale")
              #prop("boolean", "useUtc")
                  #options()
                      #item_desc ("true", "the hardware clock uses UTC")
                      #item_desc ("false", "the hardware clock does not use UTC")
                  #options_end()
          #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getLocale: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
