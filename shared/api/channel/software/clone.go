package software

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Clone a channel.  If arch_label is omitted, the arch label of the
      original channel will be used. If parent_label is omitted, the clone will be
      a base channel.
func Clone(cnxDetails *api.ConnectionDetails, OriginalLabel string, OriginalState bool) (*types.#param_desc("int", "id", "the cloned channel ID"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"originalLabel":       OriginalLabel,
		"originalState":       OriginalState,
	}

	res, err := api.Post[types.#param_desc("int", "id", "the cloned channel ID")](client, "channel/software", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute clone: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
