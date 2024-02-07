package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Update the package state of a given system
                          (High state would be needed to actually install/remove the package)
func UpdatePackageState(cnxDetails *api.ConnectionDetails, Sid int, PackageName string, State int, VersionConstraint int) (*types.1 on success, exception on failure, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"packageName":       PackageName,
		"state":       State,
		"versionConstraint":       VersionConstraint,
	}

	res, err := api.Post[types.1 on success, exception on failure](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute updatePackageState: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
