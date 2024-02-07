package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Creates a software channel
func Create(cnxDetails *api.ConnectionDetails, Label string, Name string, Summary string, ArchLabel string, ParentLabel string, $param.getFlagName() $param.getType(), GpgCheck bool) (*types.#param_desc("int", "status", "1 if the creation operation succeeded, 0 otherwise"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":       Label,
		"name":       Name,
		"summary":       Summary,
		"archLabel":       ArchLabel,
		"parentLabel":       ParentLabel,
		"$param.getName()":       $param.getFlagName(),
		"gpgCheck":       GpgCheck,
	}

	res, err := api.Post[types.#param_desc("int", "status", "1 if the creation operation succeeded, 0 otherwise")](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
