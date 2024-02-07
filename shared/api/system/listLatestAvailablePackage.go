package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the latest available version of a package for each system
func ListLatestAvailablePackage(cnxDetails *api.ConnectionDetails, Sids []int, PackageName string) (*types.#return_array_begin()
         #struct_begin("system")
             #prop_desc("int", "id", "server ID")
             #prop_desc("string", "name", "server name")
             #prop_desc("struct", "package", "package structure")
                 #struct_begin("package")
                     #prop("int", "id")
                     #prop("string", "name")
                     #prop("string", "version")
                     #prop("string", "release")
                     #prop("string", "epoch")
                     #prop("string", "arch")
                #struct_end()
        #struct_end()
    #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sids {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if PackageName {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
         #struct_begin("system")
             #prop_desc("int", "id", "server ID")
             #prop_desc("string", "name", "server name")
             #prop_desc("struct", "package", "package structure")
                 #struct_begin("package")
                     #prop("int", "id")
                     #prop("string", "name")
                     #prop("string", "version")
                     #prop("string", "release")
                     #prop("string", "epoch")
                     #prop("string", "arch")
                #struct_end()
        #struct_end()
    #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listLatestAvailablePackage: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
