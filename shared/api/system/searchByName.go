package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Returns a list of system IDs whose name matches
  the supplied regular expression(defined by
  http://docs.oracle.com/javase/1.5.0/docs/api/java/util/regex/Pattern.html_blank
 Java representation of regular expressions)
func SearchByName(cnxDetails *api.ConnectionDetails, Regexp string) (*types.#return_array_begin()
              $ShortSystemInfoSerializer
          #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"regexp":       Regexp,
	}

	res, err := api.Post[types.#return_array_begin()
              $ShortSystemInfoSerializer
          #array_end()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute searchByName: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
