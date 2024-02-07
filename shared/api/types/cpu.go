package types

type Cpu struct {
	Cache       string `mapstructure:"cache"`
	Family      string `mapstructure:"family"`
	Mhz         string `mapstructure:"mhz"`
	Flags       string `mapstructure:"flags"`
	Model       string `mapstructure:"model"`
	Vendor      string `mapstructure:"vendor"`
	Arch        string `mapstructure:"arch"`
	Stepping    string `mapstructure:"stepping"`
	Count       string `mapstructure:"count"`
	SocketCount int    `mapstructure:"socket_count"`
	CoreCount   int    `mapstructure:"core_count"`
	ThreadCount int    `mapstructure:"thread_count"`
}
