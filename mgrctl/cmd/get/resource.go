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
	Filter     string
	PageNumber int
	PageSize   int
	SortBy     string

	OutputFormat string
	Out          io.Writer // rendered items go here; tests pass a buffer
}

// ResourceFetcher is the strongly typed interface to implement for each resource.
type ResourceFetcher[R any] interface {
	List(client *api.APIClient, opts ListOptions) ([]R, int, error)
	Columns() []utils.ColumnDef
	Help() ResourceHelp
}

// ResourceHelp describes a resource in command help and resource listings.
// Details and Examples are optional so simple resources only need a summary.
type ResourceHelp struct {
	Summary  string
	Details  string
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
func registerResource[R any](name string, fetcher ResourceFetcher[R], aliases []string) {
	resourceTypes[name] = resource{
		Columns: fetcher.Columns,
		Help:    fetcher.Help,
		Aliases: aliases,
		ListAndPrint: func(client *api.APIClient, opts ListOptions) error {
			items, total, err := fetcher.List(client, opts)
			if err != nil {
				return err
			}

			if total > 0 && opts.PageSize > 0 {
				log.Info().Msgf(L("Fetched %[1]d items out of %[2]d total"), len(items), total)
			}
			return utils.PrintOutput(opts.OutputFormat, items, fetcher.Columns(), opts.Out)
		},
	}
}

// registeredTypes returns all valid resource names and aliases for cobra argument validation.
func registeredTypes() []string {
	names := make([]string, 0)
	for name, resource := range resourceTypes {
		names = append(names, name)
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

func getResourceHelpDetails() string {
	var details []string
	for _, resource := range GetRegisteredResources() {
		if resource.Help.Details == "" {
			continue
		}
		details = append(details, fmt.Sprintf(
			"  %s\n    %s",
			resource.Name,
			strings.ReplaceAll(resource.Help.Details, "\n", "\n    "),
		))
	}
	return strings.Join(details, "\n\n")
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
