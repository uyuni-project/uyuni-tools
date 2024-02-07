package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Disassociates a repository from a channel
func DisassociateRepo(cnxDetails *api.ConnectionDetails, ChannelLabel string, RepoLabel string) (*types.Channel, error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"channelLabel":       ChannelLabel,
		"repoLabel":       RepoLabel,
	}

	res, err := api.Post[types.Channel](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute disassociateRepo: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
