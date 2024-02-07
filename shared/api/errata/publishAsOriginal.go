package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Adds an existing cloned errata to a set of cloned
 channels according to its original erratum
func PublishAsOriginal(cnxDetails *api.ConnectionDetails, AdvisoryName string, $param.getFlagName() $param.getType()) (*types.Errata, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"advisoryName":       AdvisoryName,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.Errata](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute publishAsOriginal: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
