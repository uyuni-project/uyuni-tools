package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Bootstrap a system for management via either Salt or Salt SSH.
 Use SSH private key for authentication.
func BootstrapWithPrivateSshKey(cnxDetails *api.ConnectionDetails, Host string, SshPort int, SshUser string, SshPrivKey string, SshPrivKeyPass string, ActivationKey string, SaltSSH bool, ProxyId int, ReactivationKey string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"host":       Host,
		"sshPort":       SshPort,
		"sshUser":       SshUser,
		"sshPrivKey":       SshPrivKey,
		"sshPrivKeyPass":       SshPrivKeyPass,
		"activationKey":       ActivationKey,
		"saltSSH":       SaltSSH,
		"proxyId":       ProxyId,
		"reactivationKey":       ReactivationKey,
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute bootstrapWithPrivateSshKey: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
