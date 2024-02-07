package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the partitioning scheme for a kickstart profile.
func GetPartitioningScheme(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#array_single("string", "a list of partitioning commands used to
 setup the partitions, logical volumes and volume groups"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "kickstart/profile/system"
	params := ""
	if KsLabel {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#array_single("string", "a list of partitioning commands used to
 setup the partitions, logical volumes and volume groups")](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getPartitioningScheme: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
