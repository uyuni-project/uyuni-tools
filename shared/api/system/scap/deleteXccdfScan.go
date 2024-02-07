package scap

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Delete OpenSCAP XCCDF Scan from the #product() database. Note that
 only those SCAP Scans can be deleted which have passed their retention period.
func DeleteXccdfScan(cnxDetails *api.ConnectionDetails, Xid int) (*types.#param_desc("boolean", "status", "indicates success of the operation"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"xid":       Xid,
	}

	res, err := api.Post[types.#param_desc("boolean", "status", "indicates success of the operation")](client, "system/scap", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deleteXccdfScan: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
