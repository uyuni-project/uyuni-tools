











package types

type NetworkInterface struct {
    Ip string `mapstructure:"ip"`
    Interface string `mapstructure:"interface"`
    Netmask string `mapstructure:"netmask"`
    HardwareAddress string `mapstructure:"hardware_address"`
    Module string `mapstructure:"module"`
    Broadcast string `mapstructure:"broadcast"`
	Ipv6 struct {
    Address string `mapstructure:"address"`
    Netmask string `mapstructure:"netmask"`
    Scope string `mapstructure:"scope"`
} `mapstructure:"ipv6"`
	Ipv4 struct {
    Address string `mapstructure:"address"`
    Netmask string `mapstructure:"netmask"`
    Broadcast string `mapstructure:"broadcast"`
} `mapstructure:"ipv4"`
} 