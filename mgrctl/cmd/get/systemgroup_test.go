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

// ids converts a slice of SystemGroups into a slice of IDs to assert more easily.
func ids(systems []SystemGroup) []int {
	res := []int{}
	for _, group := range systems {
		res = append(res, group.ID)
	}
	return res
}

func TestSystemGroupListResponse(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
			testutils.AssertEquals(t, "wrong system group endpoint", "/rhn/manager/api/systemgroup/listAllGroups", req.URL.Path)
			return testutils.GetResponse(200, `{
				"success": true,
				"result": [
					{"id": 1, "name": "zulu"},
					{"id": 2, "name": "alpha"}
				]
			}`)
		}},
	}

	groups, total, err := (&systemGroupFetcher{}).List(client, ListOptions{SortBy: "name"})
	testutils.AssertNoError(t, "unexpected error", err)
	testutils.AssertEquals(t, "wrong number of returned items", 2, total)
	testutils.AssertEquals(t, "wrong returned IDs", []int{1, 2}, ids(groups))
}

func TestSystemGroupListErrorResponse(t *testing.T) {
	client := &api.APIClient{
		BaseURL: "https://example.test/rhn/manager/api",
		Client: &mocks.MockClient{DoFunc: func(req *http.Request) (*http.Response, error) {
			testutils.AssertEquals(t, "wrong system group endpoint", "/rhn/manager/api/systemgroup/listAllGroups", req.URL.Path)
			return testutils.GetResponse(200, `{
				"success": false,
				"message": "Some error from the server"
			}`)
		}},
	}

	_, _, err := (&systemGroupFetcher{}).List(client, ListOptions{SortBy: "name"})
	testutils.AssertError(t, "Some error from the server", err)
}
