package trusts

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove an organization to the list of trusted organizations.
func RemoveTrust(cnxDetails *api.ConnectionDetails, OrgId int, TrustOrgId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"orgId":       OrgId,
		"trustOrgId":       TrustOrgId,
	}

	res, err := api.Post[types.#return_int_success()](client, "org/trusts", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeTrust: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
