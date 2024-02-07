package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Merges all errata from one channel into another
func MergeErrata(cnxDetails *api.ConnectionDetails, MergeFromLabel string, MergeToLabel string, StartDate string, EndDate string, $param.getFlagName() $param.getType()) (*types.#return_array_begin()
          $ErrataSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"mergeFromLabel":       MergeFromLabel,
		"mergeToLabel":       MergeToLabel,
		"startDate":       StartDate,
		"endDate":       EndDate,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_array_begin()
          $ErrataSerializer
      #array_end()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute mergeErrata: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
