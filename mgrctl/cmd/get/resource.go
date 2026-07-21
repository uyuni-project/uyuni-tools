// SPDX-FileCopyrightText: 2026 Jayprakash
//
// SPDX-License-Identifier: Apache-2.0

package get

import (
	"fmt"
	"sort"
	"strings"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type ResourceFetcher interface {
	List(client *api.APIClient, filter string, page, pageSize int) ([]map[string]any, int, error)
	Columns() []utils.ColumnDef
}

type Resource struct {
	Fetcher     ResourceFetcher
	Aliases     []string
	Description string
}

var resourceTypes = make(map[string]Resource)

// registerResource adds a resource type to the global lookup table.
// Called from each resource file's init() so that adding a resource
// only requires creating a new file with no changes to resource.go.
func registerResource(name string, fetcher ResourceFetcher, aliases []string, description string) {
	resourceTypes[name] = Resource{
		Fetcher:     fetcher,
		Aliases:     aliases,
		Description: description,
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

// ResourceInfo holds exported metadata about a registered resource.
type ResourceInfo struct {
	Name        string
	Aliases     []string
	Description string
}

// GetRegisteredResources returns a sorted list of all registered resources.
func GetRegisteredResources() []ResourceInfo {
	var resources []ResourceInfo
	for name, res := range resourceTypes {
		resources = append(resources, ResourceInfo{
			Name:        name,
			Aliases:     res.Aliases,
			Description: res.Description,
		})
	}
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})
	return resources
}

func lookupFetcher(name string) (ResourceFetcher, error) {
	if res, ok := resourceTypes[name]; ok {
		return res.Fetcher, nil
	}
	for _, res := range resourceTypes {
		if utils.Contains(res.Aliases, name) {
			return res.Fetcher, nil
		}
	}
	return nil, fmt.Errorf(L("unknown resource type %[1]q; available: %[2]s"), name, registeredTypesText())
}
