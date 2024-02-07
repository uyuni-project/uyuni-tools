package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Merges all packages from one channel into another
func MergePackages(cnxDetails *api.ConnectionDetails, MergeFromLabel string, MergeToLabel string, AlignModules bool) (*types.#return_array_begin()
          $PackageSerializer
      #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"mergeFromLabel":       MergeFromLabel,
		"mergeToLabel":       MergeToLabel,
		"alignModules":       AlignModules,
	}

	res, err := api.Post[types.#return_array_begin()
          $PackageSerializer
      #array_end()](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute mergePackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
