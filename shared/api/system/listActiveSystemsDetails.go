package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Given a list of server ids, returns a list of active servers'
 details visible to the user.
func ListActiveSystemsDetails(cnxDetails *api.ConnectionDetails, Sids []int) (*types.#return_array_begin()
     #struct_begin("server details")
       #prop_desc("int", "id", "The server's id")
       #prop_desc("string", "name", "The server's name")
       #prop_desc("boolean", "payg", "Whether the server instance is payg or not")
       #prop_desc("string", "minion_id", "The server's minion id, in case it is a salt minion client")
       #prop_desc("$date", "last_checkin",
         "Last time server successfully checked in (in UTC)")
       #prop_desc("int", "ram", "The amount of physical memory in MB.")
       #prop_desc("int", "swap", "The amount of swap space in MB.")
       #prop_desc("struct", "network_devices", "The server's network devices")
       $NetworkInterfaceSerializer
       #prop_desc("struct", "dmi_info", "The server's dmi info")
       $DmiSerializer
       #prop_desc("struct", "cpu_info", "The server's cpu info")
       $CpuSerializer
       #prop_desc("array", "subscribed_channels", "List of subscribed channels")
         #return_array_begin()
           #struct_begin("channel")
             #prop_desc("int", "channel_id", "The channel id.")
             #prop_desc("string", "channel_label", "The channel label.")
           #struct_end()
         #array_end()
       #prop_desc("array", "active_guest_system_ids",
           "List of virtual guest system ids for active guests")
         #return_array_begin()
           #prop_desc("int", "guest_id", "The guest's system id.")
         #array_end()
     #struct_end()
   #array_end(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	query := "system"
	params := ""
	if Sids {
		params := fmt.Sprintf("%s&%s=%s", params, "$param.getName", $param.getFlagName)
	}
	if params != "" {
		query := fmt.Sprintf("%s?%s", query, params)
	}

    res, err := api.Get[types.#return_array_begin()
     #struct_begin("server details")
       #prop_desc("int", "id", "The server's id")
       #prop_desc("string", "name", "The server's name")
       #prop_desc("boolean", "payg", "Whether the server instance is payg or not")
       #prop_desc("string", "minion_id", "The server's minion id, in case it is a salt minion client")
       #prop_desc("$date", "last_checkin",
         "Last time server successfully checked in (in UTC)")
       #prop_desc("int", "ram", "The amount of physical memory in MB.")
       #prop_desc("int", "swap", "The amount of swap space in MB.")
       #prop_desc("struct", "network_devices", "The server's network devices")
       $NetworkInterfaceSerializer
       #prop_desc("struct", "dmi_info", "The server's dmi info")
       $DmiSerializer
       #prop_desc("struct", "cpu_info", "The server's cpu info")
       $CpuSerializer
       #prop_desc("array", "subscribed_channels", "List of subscribed channels")
         #return_array_begin()
           #struct_begin("channel")
             #prop_desc("int", "channel_id", "The channel id.")
             #prop_desc("string", "channel_label", "The channel label.")
           #struct_end()
         #array_end()
       #prop_desc("array", "active_guest_system_ids",
           "List of virtual guest system ids for active guests")
         #return_array_begin()
           #prop_desc("int", "guest_id", "The guest's system id.")
         #array_end()
     #struct_end()
   #array_end()](client, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute listActiveSystemsDetails: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
