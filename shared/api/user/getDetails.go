package user

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns the details about a given user.
func GetDetails(cnxDetails *api.ConnectionDetails, Login string) (*types.#struct_begin("user details")
     #prop_desc("string", "first_names", "deprecated, use first_name")
     #prop("string", "first_name")
     #prop("string", "last_name")
     #prop("string", "email")
     #prop("int", "org_id")
     #prop("string", "org_name")
     #prop("string", "prefix")
     #prop("string", "last_login_date")
     #prop("string", "created_date")
     #prop_desc("boolean", "enabled", "true if user is enabled,
     false if the user is disabled")
     #prop_desc("boolean", "use_pam", "true if user is configured to use
     PAM authentication")
     #prop_desc("boolean", "read_only", "true if user is readonly")
     #prop_desc("boolean", "errata_notification", "true if errata e-mail notification
     is enabled for the user")
   #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "user"
	params := ""
	if Login {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("user details")
     #prop_desc("string", "first_names", "deprecated, use first_name")
     #prop("string", "first_name")
     #prop("string", "last_name")
     #prop("string", "email")
     #prop("int", "org_id")
     #prop("string", "org_name")
     #prop("string", "prefix")
     #prop("string", "last_login_date")
     #prop("string", "created_date")
     #prop_desc("boolean", "enabled", "true if user is enabled,
     false if the user is disabled")
     #prop_desc("boolean", "use_pam", "true if user is configured to use
     PAM authentication")
     #prop_desc("boolean", "read_only", "true if user is readonly")
     #prop_desc("boolean", "errata_notification", "true if errata e-mail notification
     is enabled for the user")
   #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
