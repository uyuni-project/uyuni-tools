package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Remove an ip range from a kickstart profile.
func RemoveIpRange(cnxDetails *api.ConnectionDetails, KsLabel string, IpAddress string) (*types.#param_desc("int", "status", "1 on successful removal, 0 if range wasn't found
 for the specified kickstart, exception otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
		"ipAddress":       IpAddress,
	}

	res, err := api.Post[types.#param_desc("int", "status", "1 on successful removal, 0 if range wasn't found
 for the specified kickstart, exception otherwise")](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute removeIpRange: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
