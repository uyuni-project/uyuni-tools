package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Externally authenticated users may be members of external groups. You
 can use these groups to assign additional roles to the users when they log in.
 Can only be called by a #product() Administrator.
func CreateExternalGroupToRoleMap(cnxDetails *api.ConnectionDetails, Name string) (*types.UserExtGroup, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.UserExtGroup](client, "user/external", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createExternalGroupToRoleMap: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
