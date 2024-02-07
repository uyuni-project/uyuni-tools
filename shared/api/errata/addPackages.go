package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add a set of packages to an erratum with the given advisory name.
 This method will only allow for modification of custom errata created either through the UI or API.
func AddPackages(cnxDetails *api.ConnectionDetails, AdvisoryName string, PackageIds []int) (*types.#param("int", "the number of packages added, exception otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"advisoryName":       AdvisoryName,
		"packageIds":       PackageIds,
	}

	res, err := api.Post[types.#param("int", "the number of packages added, exception otherwise")](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addPackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
