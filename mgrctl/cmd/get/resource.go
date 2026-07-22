// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type Resource interface {
}

// ResourceFetcher is as strongly typed interface to implement for each resource.
type ResourceFetcher[R any] interface {
	List(client *api.APIClient, filter string, page int, pageSize int) ([]R, int, error)
	Columns() []utils.ColumnDef
}

type resource struct {
	ListAndPrint func(client *api.APIClient, filter string, page int, pageSize int, outputFormat string, out io.Writer) error
	Columns      func() []utils.ColumnDef
	Aliases      []string
	Description  string
}

var resourceTypes = make(map[string]resource)

// registerResource adds a resource type to the global lookup table.
// Called from each resource file's init() so that adding a resource
// only requires creating a new file with no changes to resource.go.
func registerResource[R any](name string, fetcher ResourceFetcher[R], aliases []string, description string) {
	resourceTypes[name] = resource{
		Columns:     fetcher.Columns,
		Aliases:     aliases,
		Description: description,
		ListAndPrint: func(
			client *api.APIClient, filter string, page int, pageSize int,
			outputFormat string, out io.Writer,
		) error {
			items, total, err := fetcher.List(client, filter, page, pageSize)
			if err != nil {
				return err
			}

			if total > 0 && pageSize > 0 {
				log.Info().Msgf("Fetched %d items out of %d total", len(items), total)
			}
			return utils.PrintOutput(outputFormat, items, fetcher.Columns(), out)
		},
	}
}

// registeredTypes returns all valid resource names and aliases for cobra argument validation.
func registeredTypes() []string {
	names := make([]string, 0)
	for name, res := range resourceTypes {
		names = append(names, name)
		names = append(names, res.Aliases...)
	}
	sort.Strings(names)
	return names
}

func registeredTypesText() string {
	return strings.Join(registeredTypes(), ", ")
}

// GetResourceHelp dynamically generates the list of available resources so we don't have to
// manually update the help text every time a new resource is added.
func GetResourceHelp() string {
	var lines []string
	for name, res := range resourceTypes {
		aliases := ""
		if len(res.Aliases) > 0 {
			aliases = fmt.Sprintf(" (%s)", strings.Join(res.Aliases, ", "))
		}
		lines = append(lines, fmt.Sprintf("  %-20s - %s", name+aliases, res.Description))
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func lookupResource(name string) (resource, error) {
	if res, ok := resourceTypes[name]; ok {
		return res, nil
	}
	for _, res := range resourceTypes {
		if utils.Contains(res.Aliases, name) {
			return res, nil
		}
	}
	return resource{}, fmt.Errorf(L("unknown resource type %[1]q; available: %[2]s"), name, registeredTypesText())
}
