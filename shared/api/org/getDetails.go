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
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Get details of organization based on organization name.
func GetOrganizationDetails(cnxDetails *api.ConnectionDetails, orgName string) (*types.Organization, error) {
	client, err := api.Init(cnxDetails)
	if err == nil {
		err = client.Login()
	}
	if err != nil {
		return nil, utils.Errorf(err, L("failed to connect to the server"))
	}
	res, err := api.Get[types.Organization](client, fmt.Sprintf("org/getDetails?name=%s", orgName))
	if err != nil {
		return nil, utils.Errorf(err, L("failed to get organization details"))
	}
	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
