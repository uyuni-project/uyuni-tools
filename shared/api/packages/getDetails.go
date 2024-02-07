package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Retrieve details for the package with the ID.
func GetDetails(cnxDetails *api.ConnectionDetails, Pid int) (*types.#struct_begin("package")
       #prop("int", "id")
       #prop("string", "name")
       #prop("string", "epoch")
       #prop("string", "version")
       #prop("string", "release")
       #prop("string", "arch_label")
       #prop_array("providing_channels", "string",
          "Channel label providing this package.")
       #prop("string", "build_host")
       #prop("string", "description")
       #prop("string", "checksum")
       #prop("string", "checksum_type")
       #prop("string", "vendor")
       #prop("string", "summary")
       #prop("string", "cookie")
       #prop("string", "license")
       #prop("string", "file")
       #prop("string", "build_date")
       #prop("string", "last_modified_date")
       #prop("string", "size")
       #prop_desc("string", "path", "The path on the #product() server's file system that
              the package resides.")
       #prop("string", "payload_size")
    #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "packages"
	params := ""
	if Pid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("package")
       #prop("int", "id")
       #prop("string", "name")
       #prop("string", "epoch")
       #prop("string", "version")
       #prop("string", "release")
       #prop("string", "arch_label")
       #prop_array("providing_channels", "string",
          "Channel label providing this package.")
       #prop("string", "build_host")
       #prop("string", "description")
       #prop("string", "checksum")
       #prop("string", "checksum_type")
       #prop("string", "vendor")
       #prop("string", "summary")
       #prop("string", "cookie")
       #prop("string", "license")
       #prop("string", "file")
       #prop("string", "build_date")
       #prop("string", "last_modified_date")
       #prop("string", "size")
       #prop_desc("string", "path", "The path on the #product() server's file system that
              the package resides.")
       #prop("string", "payload_size")
    #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
