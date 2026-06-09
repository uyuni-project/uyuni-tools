// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func newClient(conn api.ConnectionDetails) (*api.APIClient, error) {
	client, err := api.Init(&conn)
	if err != nil {
		return nil, err
	}
	if client.Details.User != "" || client.Details.InSession {
		err = client.Login()
	}
	if err != nil {
		return nil, utils.Errorf(err, L("unable to login to the server"))
	}
	return client, nil
}
