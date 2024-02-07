package packages

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List the dependencies for a package.
func ListDependencies(cnxDetails *api.ConnectionDetails, Pid int) (*types.#return_array_begin()
     #struct_begin("dependency")
       #prop("string", "dependency")
       #prop_desc("string", "dependency_type", "One of the following:")
         #options()
           #item("requires")
           #item("conflicts")
           #item("obsoletes")
           #item("provides")
           #item("recommends")
           #item("suggests")
           #item("supplements")
           #item("enhances")
           #item("predepends")
           #item("breaks")
         #options_end()
       #prop("string", "dependency_modifier")
     #struct_end()
   #array_end(), error) {
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

    res, err := api.Get[types.#return_array_begin()
     #struct_begin("dependency")
       #prop("string", "dependency")
       #prop_desc("string", "dependency_type", "One of the following:")
         #options()
           #item("requires")
           #item("conflicts")
           #item("obsoletes")
           #item("provides")
           #item("recommends")
           #item("suggests")
           #item("supplements")
           #item("enhances")
           #item("predepends")
           #item("breaks")
         #options_end()
       #prop("string", "dependency_modifier")
     #struct_end()
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listDependencies: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
