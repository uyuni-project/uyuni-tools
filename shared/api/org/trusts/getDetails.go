package trusts

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// The trust details about an organization given
 the organization's ID.
func GetDetails(cnxDetails *api.ConnectionDetails, OrgId int) (*types.#struct_begin("org trust details")
          #prop_desc("$date", "created", "Date the organization was
          created")
          #prop_desc("$date", "trusted_since", "Date the organization was
          defined as trusted")
          #prop_desc("int", "channels_provided", "Number of channels provided by
          the organization.")
          #prop_desc("int", "channels_consumed", "Number of channels consumed by
          the organization.")
          #prop_desc("int", "systems_migrated_to", "(Deprecated by systems_transferred_to) Number
          of systems transferred to the organization.")
          #prop_desc("int", "systems_migrated_from", "(Deprecated by systems_transferred_from) Number
          of systems transferred from the organization.")
          #prop_desc("int", "systems_transferred_to", "Number of systems transferred to
          the organization.")
          #prop_desc("int", "systems_transferred_from", "Number of systems transferred
          from the organization.")
     #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "org/trusts"
	params := ""
	if OrgId {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("org trust details")
          #prop_desc("$date", "created", "Date the organization was
          created")
          #prop_desc("$date", "trusted_since", "Date the organization was
          defined as trusted")
          #prop_desc("int", "channels_provided", "Number of channels provided by
          the organization.")
          #prop_desc("int", "channels_consumed", "Number of channels consumed by
          the organization.")
          #prop_desc("int", "systems_migrated_to", "(Deprecated by systems_transferred_to) Number
          of systems transferred to the organization.")
          #prop_desc("int", "systems_migrated_from", "(Deprecated by systems_transferred_from) Number
          of systems transferred from the organization.")
          #prop_desc("int", "systems_transferred_to", "Number of systems transferred to
          the organization.")
          #prop_desc("int", "systems_transferred_from", "Number of systems transferred
          from the organization.")
     #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
