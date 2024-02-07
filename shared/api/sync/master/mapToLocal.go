package master

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add a single organizations to the list of those the specified Master has
 exported to this Slave
func MapToLocal(cnxDetails *api.ConnectionDetails, MasterId int, MasterOrgId int, LocalOrgId int) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"masterId":       MasterId,
		"masterOrgId":       MasterOrgId,
		"localOrgId":       LocalOrgId,
	}

	res, err := api.Post[types.#return_int_success()](client, "sync/master", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute mapToLocal: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
