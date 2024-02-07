package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of the packages affected by the errata with the given advisory name.
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the packages of both of them.
func ListPackages(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#return_array_begin()
              #struct_begin("package")
                  #prop("int", "id")
                  #prop("string", "name")
                  #prop("string", "epoch")
                  #prop("string", "version")
                  #prop("string", "release")
                  #prop("string", "arch_label")
                  #prop_array("providing_channels", "string", "- Channel label
                              providing this package.")
                  #prop("string", "build_host")
                  #prop("string", "description")
                  #prop("string", "checksum")
                  #prop("string", "checksum_type")
                  #prop("string", "vendor")
                  #prop("string", "summary")
                  #prop("string", "cookie")
                  #prop("string", "license")
                  #prop("string", "path")
                  #prop("string", "file")
                  #prop("string", "build_date")
                  #prop("string", "last_modified_date")
                  #prop("string", "size")
                  #prop("string", "payload_size")
               #struct_end()
           #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "errata"
	params := ""
	if AdvisoryName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
              #struct_begin("package")
                  #prop("int", "id")
                  #prop("string", "name")
                  #prop("string", "epoch")
                  #prop("string", "version")
                  #prop("string", "release")
                  #prop("string", "arch_label")
                  #prop_array("providing_channels", "string", "- Channel label
                              providing this package.")
                  #prop("string", "build_host")
                  #prop("string", "description")
                  #prop("string", "checksum")
                  #prop("string", "checksum_type")
                  #prop("string", "vendor")
                  #prop("string", "summary")
                  #prop("string", "cookie")
                  #prop("string", "license")
                  #prop("string", "path")
                  #prop("string", "file")
                  #prop("string", "build_date")
                  #prop("string", "last_modified_date")
                  #prop("string", "size")
                  #prop("string", "payload_size")
               #struct_end()
           #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listPackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
