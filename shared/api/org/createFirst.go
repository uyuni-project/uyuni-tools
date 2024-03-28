// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

// Create first organization and user after initial setup without authentication.
// orgName is the name of the first organization to create and admin the user to create.
func CreateFirst(cnxDetails *api.ConnectionDetails, orgName string, admin *types.User) (*types.Organization, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf(L("failed to connect to the server: %s"), err)
	}

	data := map[string]interface{}{
		"orgName":       orgName,
		"adminLogin":    admin.Login,
		"adminPassword": admin.Password,
		"firstName":     admin.FirstName,
		"lastName":      admin.LastName,
		"email":         admin.Email,
	}

	res, err := api.Post[types.Organization](client, "org/createFirst", data)
	if err != nil {
		return nil, fmt.Errorf(L("failed to create first user and organization: %s"), err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
