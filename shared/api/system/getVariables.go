package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Lists kickstart variables set  in the system record
  for the specified server.
  Note: This call assumes that a system record exists in cobbler for the
  given system and will raise an XMLRPC fault if that is not the case.
  To create a system record over xmlrpc use system.createSystemRecord

  To create a system record in the Web UI  please go to
  System -&gt; &lt;Specified System&gt; -&gt; Provisioning -&gt;
  Select a Kickstart profile -&gt; Create Cobbler System Record.
func GetVariables(cnxDetails *api.ConnectionDetails, Sid int) (*types.#struct_begin("System kickstart variables")
          #prop_desc("boolean" "netboot" "netboot enabled")
          #prop_array_begin("kickstart variables")
              #struct_begin("kickstart variable")
                  #prop("string", "key")
                  #prop("string or int", "value")
              #struct_end()
          #prop_array_end()
      #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#struct_begin("System kickstart variables")
          #prop_desc("boolean" "netboot" "netboot enabled")
          #prop_array_begin("kickstart variables")
              #struct_begin("kickstart variable")
                  #prop("string", "key")
                  #prop("string or int", "value")
              #struct_end()
          #prop_array_end()
      #struct_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getVariables: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
