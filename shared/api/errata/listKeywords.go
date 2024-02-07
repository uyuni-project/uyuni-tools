package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the keywords associated with an erratum matching the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the keywords of both of them.
func ListKeywords(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#array_single("string", "keyword associated with erratum."), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "errata"
	params := ""
	if AdvisoryName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "keyword associated with erratum.")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listKeywords: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
