package external

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Externally authenticated users may be members of external groups. You
 can use these groups to give access to server groups to the users when they log in.
 Can only be called by an org_admin.
func CreateExternalGroupToSystemGroupMap(cnxDetails *api.ConnectionDetails, Name string) (*types.OrgUserExtGroup, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"name":       Name,
	}

	res, err := api.Post[types.OrgUserExtGroup](client, "user/external", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createExternalGroupToSystemGroupMap: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
