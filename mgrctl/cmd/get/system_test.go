// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"net/http"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/mocks"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

func TestSystemListPassesSortingToAPI(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
			testutils.AssertEquals(
				t,
				"wrong system endpoint",
				"/rhn/manager/api/system/listSystemsFiltered",
				req.URL.Path,
			)

			query := req.URL.Query()
			testutils.AssertEquals(t, "wrong sort key", "name", query.Get("sortKey"))
			testutils.AssertEquals(t, "wrong sort direction", "false", query.Get("sortDescending"))

			return testutils.GetResponse(
				http.StatusOK,
				`{"success":true,"result":{"data":[],"total":0}}`,
			)
		}},
	}

	_, _, err := (&systemFetcher{}).List(client, ListOptions{SortBy: "name"})
	testutils.AssertNoError(t, "unexpected error", err)
}
