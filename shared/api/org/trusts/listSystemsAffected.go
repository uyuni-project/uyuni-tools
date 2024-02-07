package trusts

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get a list of systems within the  trusted organization
   that would be affected if the trust relationship was removed.
   This basically lists systems that are sharing at least (1) package.
func ListSystemsAffected(cnxDetails *api.ConnectionDetails, OrgId int, TrustOrgId string) (*types.#return_array_begin()
     #struct_begin("affected systems")
       #prop("int", "systemId")
       #prop("string", "systemName")
     #struct_end()
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "org/trusts"
	params := ""
	if OrgId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if TrustOrgId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     #struct_begin("affected systems")
       #prop("int", "systemId")
       #prop("string", "systemName")
     #struct_end()
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSystemsAffected: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
