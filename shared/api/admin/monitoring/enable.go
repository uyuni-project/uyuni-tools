package monitoring

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Enable monitoring.
func Enable(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
      #struct_begin("Exporters")
          #prop("string", "node")
          #prop("string", "tomcat")
          #prop("string", "taskomatic")
          #prop("string", "postgres")
          #prop("string", "self_monitoring")
      #struct_end()
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
	}

	res, err := api.Post[types.#return_array_begin()
      #struct_begin("Exporters")
          #prop("string", "node")
          #prop("string", "tomcat")
          #prop("string", "taskomatic")
          #prop("string", "postgres")
          #prop("string", "self_monitoring")
      #struct_end()
  #array_end()](client, "admin/monitoring", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute enable: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
