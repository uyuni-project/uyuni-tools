package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Fetch the playbook content from the control node using a synchronous salt call.
func FetchPlaybookContents(cnxDetails *api.ConnectionDetails, PathId int, PlaybookRelPath string) (*types.#param_desc("string", "contents", "Text contents of the playbook"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"pathId":       PathId,
		"playbookRelPath":       PlaybookRelPath,
	}

	res, err := api.Post[types.#param_desc("string", "contents", "Text contents of the playbook")](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute fetchPlaybookContents: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
