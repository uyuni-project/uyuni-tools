package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Set advanced options for a kickstart profile.
 'md5_crypt_rootpw' is not supported anymore.
 If 'sha256_crypt_rootpw' is set to 'True', 'root_pw' is taken as plaintext and
 will sha256 encrypted on server side, otherwise a hash encoded password
 (according to the auth option) is expected
func SetAdvancedOptions(cnxDetails *api.ConnectionDetails, KsLabel string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setAdvancedOptions: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
