package slave

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Updates attributes of the specified Slave
func Update(cnxDetails *api.ConnectionDetails, SlaveId int, SlaveFqdn string, IsEnabled bool, AllowAllOrgs bool) (*types.IssSlave, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"slaveId":       SlaveId,
		"slaveFqdn":       SlaveFqdn,
		"isEnabled":       IsEnabled,
		"allowAllOrgs":       AllowAllOrgs,
	}

	res, err := api.Post[types.IssSlave](client, "sync/slave", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
