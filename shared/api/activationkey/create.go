package activationkey

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Create a new activation key.
 The activation key parameter passed
 in will be prefixed with the organization ID, and this value will be
 returned from the create call.

 Eg. If the caller passes in the key "foo" and belong to an organization with
 the ID 100, the actual activation key will be "100-foo".

 This call allows for the setting of a usage limit on this activation key.
 If unlimited usage is desired see the similarly named API method with no
 usage limit argument.
func Create(cnxDetails *api.ConnectionDetails, Key string, Description string, BaseChannelLabel string, UsageLimit int, UniversalDefault bool, Key string, $param.getFlagName() $param.getType()) (*types.#param("string", "The new activation key"), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"key":       Key,
		"description":       Description,
		"baseChannelLabel":       BaseChannelLabel,
		"usageLimit":       UsageLimit,
		"universalDefault":       UniversalDefault,
		"key":       Key,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#param("string", "The new activation key")](client, "activationkey", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute create: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
