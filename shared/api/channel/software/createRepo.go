package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Creates a repository
func CreateRepo(cnxDetails *api.ConnectionDetails, Label string, Type string, Url string, Type string, SslCaCert string, SslCliCert string, SslCliKey string, Type string, SslCaCert string, SslCliCert string, SslCliKey string, HasSignedMetadata bool) (*types.ContentSource, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"label":             Label,
		"type":              Type,
		"url":               Url,
		"type":              Type,
		"sslCaCert":         SslCaCert,
		"sslCliCert":        SslCliCert,
		"sslCliKey":         SslCliKey,
		"type":              Type,
		"sslCaCert":         SslCaCert,
		"sslCliCert":        SslCliCert,
		"sslCliKey":         SslCliKey,
		"hasSignedMetadata": HasSignedMetadata,
	}

	res, err := api.Post[types.ContentSource](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute createRepo: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
