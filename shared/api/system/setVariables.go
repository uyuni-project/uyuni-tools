package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Sets a list of kickstart variables in the cobbler system record
 for the specified server.
  Note: This call assumes that a system record exists in cobbler for the
  given system and will raise an XMLRPC fault if that is not the case.
  To create a system record over xmlrpc use system.createSystemRecord

  To create a system record in the Web UI  please go to
  System -&gt; &lt;Specified System&gt; -&gt; Provisioning -&gt;
  Select a Kickstart profile -&gt; Create Cobbler System Record.
func SetVariables(cnxDetails *api.ConnectionDetails, Sid int, Netboot bool, $param.getFlagName() $param.getType()) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"netboot":       Netboot,
		"$param.getName()":       $param.getFlagName(),
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute setVariables: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
