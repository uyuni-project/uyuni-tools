package proxy

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Compute and download the configuration for proxy containers
func ContainerConfig(cnxDetails *api.ConnectionDetails, ProxyName string, ProxyPort int, Server string, MaxCache int, Email string, RootCA string, ProxyCrt string, ProxyKey string, CaCrt string, CaKey string, CaPassword string, $param.getFlagName() $param.getType(), Country string, State string, City string, Org string, OrgUnit string, SslEmail string) (*types.#array_single("byte", "binary object - package file"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"proxyName":       ProxyName,
		"proxyPort":       ProxyPort,
		"server":       Server,
		"maxCache":       MaxCache,
		"email":       Email,
		"rootCA":       RootCA,
		"proxyCrt":       ProxyCrt,
		"proxyKey":       ProxyKey,
		"caCrt":       CaCrt,
		"caKey":       CaKey,
		"caPassword":       CaPassword,
		"$param.getName()":       $param.getFlagName(),
		"country":       Country,
		"state":       State,
		"city":       City,
		"org":       Org,
		"orgUnit":       OrgUnit,
		"sslEmail":       SslEmail,
	}

	res, err := api.Post[types.#array_single("byte", "binary object - package file")](client, "proxy", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute containerConfig: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
