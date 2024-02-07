package tree

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Edit a Kickstart Tree (Distribution) in #product().
func Update(cnxDetails *api.ConnectionDetails, TreeLabel string, BasePath string, ChannelLabel string, InstallType string, InstallType string, KernelOptions string, PostKernelOptions string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"treeLabel":       TreeLabel,
		"basePath":       BasePath,
		"channelLabel":       ChannelLabel,
		"installType":       InstallType,
		"installType":       InstallType,
		"kernelOptions":       KernelOptions,
		"postKernelOptions":       PostKernelOptions,
	}

	res, err := api.Post[types.#return_int_success()](client, "kickstart/tree", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
