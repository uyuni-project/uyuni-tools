package errata

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Get the Bugzilla fixes for an erratum matching the given
 advisoryName. The bugs will be returned in a struct where the bug id is
 the key.  i.e. 208144="errata.bugzillaFixes Method Returns different
 results than docs say"
 For those errata that are present in both vendor and user organizations under the same advisory name,
 this method retrieves the list of Bugzilla fixes of both of them.
func BugzillaFixes(cnxDetails *api.ConnectionDetails, AdvisoryName string) (*types.#struct_begin("Bugzilla info")
          #prop_desc("string", "bugzilla_id", "actual bug number is the key into the
                      struct")
          #prop_desc("string", "bug_summary", "summary who's key is the bug id")
      #struct_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"advisoryName":       AdvisoryName,
	}

	res, err := api.Post[types.#struct_begin("Bugzilla info")
          #prop_desc("string", "bugzilla_id", "actual bug number is the key into the
                      struct")
          #prop_desc("string", "bug_summary", "summary who's key is the bug id")
      #struct_end()](client, "errata", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute bugzillaFixes: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
