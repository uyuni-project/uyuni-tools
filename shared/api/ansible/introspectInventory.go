package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Introspect inventory under given inventory path with given pathId and return it in a structured way
func IntrospectInventory(cnxDetails *api.ConnectionDetails, PathId int) (*types.#struct_begin("Inventory in a nested structure")
   #param_desc("object", "Inventory item", "Inventory item (can be nested)")
 #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"pathId":       PathId,
	}

	res, err := api.Post[types.#struct_begin("Inventory in a nested structure")
   #param_desc("object", "Inventory item", "Inventory item (can be nested)")
 #struct_end()](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute introspectInventory: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
