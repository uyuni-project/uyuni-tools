package search

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Advanced method to search lucene indexes with a passed in query written
 in Lucene Query Parser syntax, additionally this method will limit results to those
 which are in the passed in channel label.
 Lucene Query Parser syntax is defined at
 http://lucene.apache.org/java/3_5_0/queryparsersyntax.html_blank
 lucene.apache.org.
 Fields searchable for Packages:
 name, epoch, version, release, arch, description, summary
 Lucene Query Example: "name:kernel AND version:2.6.18 AND -description:devel"
func AdvancedWithChannel(cnxDetails *api.ConnectionDetails, LuceneQuery string, ChannelLabel string) (*types.#return_array_begin()
      $PackageOverviewSerializer
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages/search"
	params := ""
	if LuceneQuery {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if ChannelLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      $PackageOverviewSerializer
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute advancedWithChannel: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
