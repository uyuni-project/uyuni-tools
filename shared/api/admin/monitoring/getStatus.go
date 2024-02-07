package monitoring

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the status of each Prometheus exporter.
func GetStatus(cnxDetails *api.ConnectionDetails) (*types.#return_array_begin()
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

	query := "admin/monitoring"
	params := ""
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      #struct_begin("Exporters")
          #prop("string", "node")
          #prop("string", "tomcat")
          #prop("string", "taskomatic")
          #prop("string", "postgres")
          #prop("string", "self_monitoring")
      #struct_end()
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute getStatus: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
