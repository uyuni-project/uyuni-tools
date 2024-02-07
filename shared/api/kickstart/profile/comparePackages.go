package profile

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list for each kickstart profile; each list will contain
             package names not present on the other profile.
func ComparePackages(cnxDetails *api.ConnectionDetails, KickstartLabel1 string, KickstartLabel2 string) (*types.#struct_begin("Comparison Info")
      #prop_desc("array", "kickstartLabel1", "Actual label of the first kickstart
                 profile is the key into the struct")
          #array_single("string", "package name")
      #prop_desc("array", "kickstartLabel2", "Actual label of the second kickstart
                 profile is the key into the struct")
          #array_single("string", "package name")
  #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"kickstartLabel1":       KickstartLabel1,
		"kickstartLabel2":       KickstartLabel2,
	}

	res, err := api.Post[types.#struct_begin("Comparison Info")
      #prop_desc("array", "kickstartLabel1", "Actual label of the first kickstart
                 profile is the key into the struct")
          #array_single("string", "package name")
      #prop_desc("array", "kickstartLabel2", "Actual label of the second kickstart
                 profile is the key into the struct")
          #array_single("string", "package name")
  #struct_end()](client, "kickstart/profile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute comparePackages: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
