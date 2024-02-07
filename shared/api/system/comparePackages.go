package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Compares the packages installed on two systems.
func ComparePackages(cnxDetails *api.ConnectionDetails, Sid1 int, Sid2 int) (*types.#return_array_begin()
              $PackageMetadataSerializer
          #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid1":       Sid1,
		"sid2":       Sid2,
	}

	res, err := api.Post[types.#return_array_begin()
              $PackageMetadataSerializer
          #array_end()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute comparePackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
