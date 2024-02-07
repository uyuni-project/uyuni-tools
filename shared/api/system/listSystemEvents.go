package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// List system actions of the specified type that were *scheduled* against the given server after the
 specified date. "actionType" should be exactly the string returned in the action_type field
 from the listSystemEvents(sessionKey, serverId) method. For example,
 'Package Install' or 'Initiate a kickstart for a virtual guest.'
 Note: see also system.getEventHistory method which returns a history of all events.
func ListSystemEvents(cnxDetails *api.ConnectionDetails, Sid int, ActionType string, EarliestDate $date) (*types.#return_array_begin()
      #struct_begin("action")
          #prop_desc("int", "failed_count", "Number of times action failed.")
          #prop_desc("string", "modified", "Date modified. (Deprecated by
                     modified_date)")
          #prop_desc($date, "modified_date", "Date modified.")
          #prop_desc("string", "created", "Date created. (Deprecated by
                     created_date)")
          #prop_desc($date, "created_date", "Date created.")
          #prop("string", "action_type")
          #prop_desc("int", "successful_count",
                     "Number of times action was successful.")
          #prop_desc("string", "earliest_action", "Earliest date this action
                     will occur.")
          #prop_desc("int", "archived", "If this action is archived. (1 or 0)")
          #prop_desc("string", "scheduler_user", "available only if concrete user
                     has scheduled the action")
          #prop_desc("string", "prerequisite", "Pre-requisite action. (optional)")
          #prop_desc("string", "name", "Name of this action.")
          #prop_desc("int", "id", "Id of this action.")
          #prop_desc("string", "version", "Version of action.")
          #prop_desc("string", "completion_time", "The date/time the event was
                     completed. Format -&gt;YYYY-MM-dd hh:mm:ss.ms
                     Eg -&gt;2007-06-04 13:58:13.0. (optional)
                     (Deprecated by completed_date)")
          #prop_desc($date, "completed_date", "The date/time the event was completed.
                     (optional)")
          #prop_desc("string", "pickup_time", "The date/time the action was picked
                     up. Format -&gt;YYYY-MM-dd hh:mm:ss.ms
                     Eg -&gt;2007-06-04 13:58:13.0. (optional)
                     (Deprecated by pickup_date)")
          #prop_desc($date, "pickup_date", "The date/time the action was picked up.
                     (optional)")
          #prop_desc("string", "result_msg", "The result string after the action
                     executes at the client machine. (optional)")
          #prop_array_begin_desc("additional_info", "This array contains additional
              information for the event, if available.")
              #struct_begin("info")
                  #prop_desc("string", "detail", "The detail provided depends on the
                  specific event.  For example, for a package event, this will be the
                  package name, for an errata event, this will be the advisory name
                  and synopsis, for a config file event, this will be path and
                  optional revision information...etc.")
                  #prop_desc("string", "result", "The result (if included) depends
                  on the specific event.  For example, for a package or errata event,
                  no result is included, for a config file event, the result might
                  include an error (if one occurred, such as the file was missing)
                  or in the case of a config file comparison it might include the
                  differences found.")
              #struct_end()
          #prop_array_end()
      #struct_end()
  #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sid {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if ActionType {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if EarliestDate {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
      #struct_begin("action")
          #prop_desc("int", "failed_count", "Number of times action failed.")
          #prop_desc("string", "modified", "Date modified. (Deprecated by
                     modified_date)")
          #prop_desc($date, "modified_date", "Date modified.")
          #prop_desc("string", "created", "Date created. (Deprecated by
                     created_date)")
          #prop_desc($date, "created_date", "Date created.")
          #prop("string", "action_type")
          #prop_desc("int", "successful_count",
                     "Number of times action was successful.")
          #prop_desc("string", "earliest_action", "Earliest date this action
                     will occur.")
          #prop_desc("int", "archived", "If this action is archived. (1 or 0)")
          #prop_desc("string", "scheduler_user", "available only if concrete user
                     has scheduled the action")
          #prop_desc("string", "prerequisite", "Pre-requisite action. (optional)")
          #prop_desc("string", "name", "Name of this action.")
          #prop_desc("int", "id", "Id of this action.")
          #prop_desc("string", "version", "Version of action.")
          #prop_desc("string", "completion_time", "The date/time the event was
                     completed. Format -&gt;YYYY-MM-dd hh:mm:ss.ms
                     Eg -&gt;2007-06-04 13:58:13.0. (optional)
                     (Deprecated by completed_date)")
          #prop_desc($date, "completed_date", "The date/time the event was completed.
                     (optional)")
          #prop_desc("string", "pickup_time", "The date/time the action was picked
                     up. Format -&gt;YYYY-MM-dd hh:mm:ss.ms
                     Eg -&gt;2007-06-04 13:58:13.0. (optional)
                     (Deprecated by pickup_date)")
          #prop_desc($date, "pickup_date", "The date/time the action was picked up.
                     (optional)")
          #prop_desc("string", "result_msg", "The result string after the action
                     executes at the client machine. (optional)")
          #prop_array_begin_desc("additional_info", "This array contains additional
              information for the event, if available.")
              #struct_begin("info")
                  #prop_desc("string", "detail", "The detail provided depends on the
                  specific event.  For example, for a package event, this will be the
                  package name, for an errata event, this will be the advisory name
                  and synopsis, for a config file event, this will be path and
                  optional revision information...etc.")
                  #prop_desc("string", "result", "The result (if included) depends
                  on the specific event.  For example, for a package or errata event,
                  no result is included, for a config file event, the result might
                  include an error (if one occurred, such as the file was missing)
                  or in the case of a config file comparison it might include the
                  differences found.")
              #struct_end()
          #prop_array_end()
      #struct_end()
  #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listSystemEvents: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
