package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of http://cve.mitre.org/_blankCVEs applicable to the errata
 with the given advisory name. For those errata that are present in both vendor and user organizations under the
 same advisory name, this method retrieves the list of CVEs of both of them.
func ListCves(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#array_single("string", "CVE name"), error) {
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

    res, err := api.Get[types.#array_single("string", "CVE name")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listCves: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
