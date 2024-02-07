package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Downloads the Cobbler-rendered Kickstart file.
func DownloadRenderedKickstart(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#param_desc("string", "ks", "The contents of the kickstart file"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
	}

	res, err := api.Post[types.#param_desc("string", "ks", "The contents of the kickstart file")](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute downloadRenderedKickstart: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
