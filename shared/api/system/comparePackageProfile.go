package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Compare a system's packages against a package profile.  In
 the result returned, 'this_system' represents the server provided as an input
 and 'other_system' represents the profile provided as an input.
func ComparePackageProfile(cnxDetails *api.ConnectionDetails, Sid int, ProfileLabel string) (*types.#return_array_begin()
              $PackageMetadataSerializer
          #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"profileLabel":       ProfileLabel,
	}

	res, err := api.Post[types.#return_array_begin()
              $PackageMetadataSerializer
          #array_end()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute comparePackageProfile: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
