// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"errors"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// CreateFirst creates the first organization and user after initial setup without authentication.
//
// orgName is the name of the first organization to create and admin the user to create.
func CreateFirst(cnxDetails *api.ConnectionDetails, orgName string, admin *types.User) (*types.Organization, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, utils.Errorf(err, L("unable to prepare API client"))
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
		return nil, utils.Errorf(err, L("failed to create first user and organization"))
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
