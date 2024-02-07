package types

type Device struct {
	Device      string `mapstructure:"device"`
	DeviceClass string `mapstructure:"device_class"`
	Driver      string `mapstructure:"driver"`
	Description string `mapstructure:"description"`
	Bus         string `mapstructure:"bus"`
	Pcitype     string `mapstructure:"pcitype"`
}
