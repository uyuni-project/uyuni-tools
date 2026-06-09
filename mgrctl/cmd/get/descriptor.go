// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package get

import "github.com/uyuni-project/uyuni-tools/shared/api"

type ColumnDef struct {
	Header string
	Field  string
}

type Resource[T any] interface {
	List(client *api.APIClient) ([]T, error)
	Get(client *api.APIClient, name string) (T, error)
	Columns() []ColumnDef
	FilterFields() []string
}
