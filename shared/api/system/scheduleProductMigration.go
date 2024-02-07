package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule a Product migration for a system. This call is the
 recommended and supported way of migrating a system to the next Service Pack. It will
 automatically find all mandatory product channels below a given target base channel
 and subscribe the system accordingly. Any additional optional channels can be
 subscribed by providing their labels.
func ScheduleProductMigration(cnxDetails *api.ConnectionDetails, Sid int, BaseChannelLabel string, OptionalChildChannels []string, DryRun bool, EarliestOccurrence $date, AllowVendorChange bool, TargetIdent string, EarliestOccurrence $date, TargetIdent string, RemoveProductsWithNoSuccessorAfterMigration bool) (*types.#param_desc("int", "actionId", "The action id of the scheduled action"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"baseChannelLabel":       BaseChannelLabel,
		"optionalChildChannels":       OptionalChildChannels,
		"dryRun":       DryRun,
		"earliestOccurrence":       EarliestOccurrence,
		"allowVendorChange":       AllowVendorChange,
		"targetIdent":       TargetIdent,
		"earliestOccurrence":       EarliestOccurrence,
		"targetIdent":       TargetIdent,
		"removeProductsWithNoSuccessorAfterMigration":       RemoveProductsWithNoSuccessorAfterMigration,
	}

	res, err := api.Post[types.#param_desc("int", "actionId", "The action id of the scheduled action")](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleProductMigration: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
