package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Creates a cobbler system record with the specified kickstart label
func CreateSystemRecord(cnxDetails *api.ConnectionDetails, Sid int, KsLabel string, SystemName string, KOptions string, Comment string, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"ksLabel":       KsLabel,
		"systemName":       SystemName,
		"kOptions":       KOptions,
		"comment":       Comment,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createSystemRecord: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
