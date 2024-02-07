package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Retrieve the package file associated with a package.
 (Consider using #getPackageUrlpackages.getPackageUrl
 for larger files.)
func GetPackage(cnxDetails *api.ConnectionDetails, Pid int) (*types.#array_single("byte", "binary object - package file"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages"
	params := ""
	if Pid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("byte", "binary object - package file")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getPackage: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
