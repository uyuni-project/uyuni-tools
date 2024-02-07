package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new organization and associated administrator account.
func Create(cnxDetails *api.ConnectionDetails, OrgName string, AdminLogin string, AdminPassword string, Prefix string, FirstName string, LastName string, Email string, UsePamAuth bool) (*types.OrgDto, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"orgName":       OrgName,
		"adminLogin":    AdminLogin,
		"adminPassword": AdminPassword,
		"prefix":        Prefix,
		"firstName":     FirstName,
		"lastName":      LastName,
		"email":         Email,
		"usePamAuth":    UsePamAuth,
	}

	res, err := api.Post[types.OrgDto](client, "org", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
