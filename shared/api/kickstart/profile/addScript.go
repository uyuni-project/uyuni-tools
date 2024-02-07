package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Add a pre/post script to a kickstart profile.
func AddScript(cnxDetails *api.ConnectionDetails, KsLabel string, Name string, Contents string, Interpreter string, Type string, Chroot bool, Template bool, Erroronfail bool) (*types.#param_desc("int", "id", "the id of the added script"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"ksLabel":       KsLabel,
		"name":       Name,
		"contents":       Contents,
		"interpreter":       Interpreter,
		"type":       Type,
		"chroot":       Chroot,
		"template":       Template,
		"erroronfail":       Erroronfail,
	}

	res, err := api.Post[types.#param_desc("int", "id", "the id of the added script")](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute addScript: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
