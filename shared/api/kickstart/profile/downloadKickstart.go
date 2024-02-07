package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Download the full contents of a kickstart file.
func DownloadKickstart(cnxDetails *api.ConnectionDetails, KsLabel string, Host string) (*types.#param_desc("string", "ks", "The contents of the kickstart file. Note: if
 an activation key is not associated with the kickstart file, registration
 will not occur in the generated %post section. If one is
 associated, it will be used for registration"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
		"host":       Host,
	}

	res, err := api.Post[types.#param_desc("string", "ks", "The contents of the kickstart file. Note: if
 an activation key is not associated with the kickstart file, registration
 will not occur in the generated %post section. If one is
 associated, it will be used for registration")](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute downloadKickstart: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
