package system

import (
	"errors"
	"fmt"

	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/types"
)

// Provision a guest on the host specified.  Defaults to:
 memory=512MB, vcpu=1, storage=3GB, mac_address=random.
func ProvisionVirtualGuest(cnxDetails *api.ConnectionDetails, Sid int, GuestName string, ProfileName string, ProfileName string, MemoryMb int, Vcpus int, StorageGb int, MacAddress string) (*types.#return_int_success(), error) {
	client, err := api.Init(cnxDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %s", err)
	}

	data := map[string]interface{}{
		"sid":       Sid,
		"guestName":       GuestName,
		"profileName":       ProfileName,
		"profileName":       ProfileName,
		"memoryMb":       MemoryMb,
		"vcpus":       Vcpus,
		"storageGb":       StorageGb,
		"macAddress":       MacAddress,
	}

	res, err := api.Post[types.#return_int_success()](client, "system", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute provisionVirtualGuest: %s", err)
	}

	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}
