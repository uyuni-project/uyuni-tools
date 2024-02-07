package slave

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new Slave, known to this Master.
func Create(cnxDetails *api.ConnectionDetails, SlaveFqdn string, IsEnabled bool, AllowAllOrgs bool) (*types.IssSlave, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"slaveFqdn":    SlaveFqdn,
		"isEnabled":    IsEnabled,
		"allowAllOrgs": AllowAllOrgs,
	}

	res, err := api.Post[types.IssSlave](client, "sync/slave", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
