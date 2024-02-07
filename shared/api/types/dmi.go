











package types

type Dmi struct {
    Vendor string `mapstructure:"vendor"`
    System string `mapstructure:"system"`
    Product string `mapstructure:"product"`
    Asset string `mapstructure:"asset"`
    Board string `mapstructure:"board"`
    BiosRelease string `mapstructure:"bios_release"`
    BiosVendor string `mapstructure:"bios_vendor"`
    BiosVersion string `mapstructure:"bios_version"`
} 