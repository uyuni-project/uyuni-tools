package scap

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Schedule OpenSCAP scan.
func ScheduleXccdfScan(cnxDetails *api.ConnectionDetails, Sids []int, XccdfPath string, OscapParams string, Date $date, XccdfPath string, OscapPrams string, OvalFiles string, Sid int) (*types.#param_desc("int", "id", "ID if SCAP action created"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sids":       Sids,
		"xccdfPath":       XccdfPath,
		"oscapParams":       OscapParams,
		"date":       Date,
		"xccdfPath":       XccdfPath,
		"oscapPrams":       OscapPrams,
		"ovalFiles":       OvalFiles,
		"sid":       Sid,
	}

	res, err := api.Post[types.#param_desc("int", "id", "ID if SCAP action created")](client, "system/scap", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scheduleXccdfScan: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
