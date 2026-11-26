// SPDX-FileCopyrightText: 2026 Jay Prakash katara <katarajayprakash@icloud.com>
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

// ListOptions holds everything a resource needs to fetch a listing and render it.
type ListOptions struct {
	// Resources whose API supports neither filtering nor pagination ignore these.
	FilterKey      string
	FilterOperator string
	FilterValue    string
	PageNumber     int
	PageSize       int
	SortBy         string
	SortDescending bool

	OutputFormat string
	Out          io.Writer // rendered items go here; tests pass a buffer
}

// hasFilter returns whether to apply filtering or not.
func (o ListOptions) hasFilter() bool {
	return o.FilterKey != "" && o.FilterValue != ""
}

// hasPagination returns whether results are paginated or not. A page size below one means no pagination.
func (o ListOptions) hasPagination() bool {
	return o.PageSize > 0
}

// pageCount returns the number of pages needed to hold total items, rounding up the last partial page.
func pageCount(total int, pageSize int) int {
	return (total + pageSize - 1) / pageSize
}

// ResourceFetcher is the strongly typed interface to implement for each resource.
type ResourceFetcher[R any] interface {
	List(client *api.APIClient, opts ListOptions) ([]R, int, error)
	Columns() []utils.ColumnDef
	Help() ResourceHelp
}

// ResourceHelp describes a resource in command help and resource listings.
// Examples are optional so simple resources only need a summary.
type ResourceHelp struct {
	Summary  string
	Examples string
}

type resource struct {
	ListAndPrint func(client *api.APIClient, opts ListOptions) error
	Columns      func() []utils.ColumnDef
	Help         func() ResourceHelp
	Aliases      []string
}

var resourceTypes = make(map[string]resource)

// registerResource adds a resource type to the global lookup table.
// Called from each resource file's init() so that adding a resource
// only requires creating a new file with no changes to resource.go.
//
// If canSort is set to true, the sort will need to happen in the List() of the fetcher,
// otherwise, it will be performed on the fetched data. The same applies to canFilter.
func registerResource[R any](name string, fetcher ResourceFetcher[R], aliases []string, canSort bool, canFilter bool) {
	resourceTypes[name] = resource{
		Columns: fetcher.Columns,
		Help:    fetcher.Help,
		Aliases: aliases,
		ListAndPrint: func(client *api.APIClient, opts ListOptions) error {
			items, total, err := fetcher.List(client, opts)
			if err != nil {
				return err
			}

			if !canSort || !canFilter {
				processor := NewDataProcessor(items, opts)
				if !canSort {
					processor = processor.Sort()
				}
				if !canFilter {
					processor = processor.Filter()
				}
				items, err = processor.Get()
				if err != nil {
					return err
				}
			}

			if total > 0 && opts.hasPagination() {
				pages := pageCount(total, opts.PageSize)
				log.Info().Msgf(L("Fetched %[1]d items out of %[2]d total. Page %[3]d/%[4]d"),
					len(items), total, opts.PageNumber+1, pages) // User wants to see page numbers starting at 1
			}
			return utils.PrintOutput(opts.OutputFormat, items, fetcher.Columns(), opts.Out)
		},
	}
}

// registeredTypes returns all valid resource names and aliases for cobra argument validation.
func registeredTypes() []string {
	names := make([]string, 0)
	for name := range resourceTypes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// registeredAliases returns all valid resource aliases for cobra argument validation.
func registeredAliases() []string {
	names := make([]string, 0)
	for _, resource := range resourceTypes {
		names = append(names, resource.Aliases...)
	}
	sort.Strings(names)
	return names
}

// ResourceInfo contains the public metadata for a registered resource.
type ResourceInfo struct {
	Name    string
	Aliases []string
	Help    ResourceHelp
}

// GetRegisteredResources returns all registered resources sorted by name.
func GetRegisteredResources() []ResourceInfo {
	resources := make([]ResourceInfo, 0, len(resourceTypes))
	for name, resource := range resourceTypes {
		resources = append(resources, ResourceInfo{
			Name:    name,
			Aliases: resource.Aliases,
			Help:    resource.Help(),
		})
	}
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})
	return resources
}

func formatResourceHeading(resource ResourceInfo) string {
	heading := resource.Name
	if len(resource.Aliases) > 0 {
		heading = fmt.Sprintf(
			"%s (%s: %s)",
			heading,
			NL("alias", "aliases", len(resource.Aliases)),
			strings.Join(resource.Aliases, ", "),
		)
	}
	return heading
}

func getResourceHelpSummaries() string {
	resources := GetRegisteredResources()
	summaries := make([]string, 0, len(resources))
	for _, resource := range resources {
		summaries = append(summaries, fmt.Sprintf("  %s\n    %s", formatResourceHeading(resource), resource.Help.Summary))
	}
	return strings.Join(summaries, "\n")
}

func getResourceHelpExamples() string {
	var examples []string
	for _, resource := range GetRegisteredResources() {
		if resource.Help.Examples != "" {
			examples = append(examples, resource.Help.Examples)
		}
	}
	return strings.Join(examples, "\n\n")
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
	return resource{}, fmt.Errorf(
		L("unknown resource type %[1]q; see 'mgrctl get --help' for available resource types"),
		name,
	)
}
