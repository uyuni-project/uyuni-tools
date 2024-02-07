











package types

type SystemSearchResult struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    LastCheckin string `mapstructure:"last_checkin"`
    Hostname string `mapstructure:"hostname"`
    Uuid string `mapstructure:"uuid"`
    Ip string `mapstructure:"ip"`
    HwDescription string `mapstructure:"hw_description"`
    HwDeviceId string `mapstructure:"hw_device_id"`
    HwVendorId string `mapstructure:"hw_vendor_id"`
    HwDriver string `mapstructure:"hw_driver"`
} 