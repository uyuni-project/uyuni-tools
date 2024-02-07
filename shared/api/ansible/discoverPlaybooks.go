package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Discover playbooks under given playbook path with given pathId
func DiscoverPlaybooks(cnxDetails *api.ConnectionDetails, PathId int) (*types.#struct_begin("playbooks")
     #struct_begin("playbook")
         $AnsiblePathSerializer
     #struct_end()
 #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"pathId":       PathId,
	}

	res, err := api.Post[types.#struct_begin("playbooks")
     #struct_begin("playbook")
         $AnsiblePathSerializer
     #struct_end()
 #struct_end()](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute discoverPlaybooks: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
