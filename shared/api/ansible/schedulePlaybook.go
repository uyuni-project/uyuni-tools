package ansible

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule a playbook execution
func SchedulePlaybook(cnxDetails *api.ConnectionDetails, PlaybookPath string, InventoryPath string, ControlNodeId int, EarliestOccurrence $date, ActionChainLabel string, TestMode bool, $param.getFlagName() $param.getType()) (*types.#param_desc("int", "id", "ID of the playbook execution action created"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"playbookPath":       PlaybookPath,
		"inventoryPath":       InventoryPath,
		"controlNodeId":       ControlNodeId,
		"earliestOccurrence":       EarliestOccurrence,
		"actionChainLabel":       ActionChainLabel,
		"testMode":       TestMode,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#param_desc("int", "id", "ID of the playbook execution action created")](client, "ansible", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schedulePlaybook: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
