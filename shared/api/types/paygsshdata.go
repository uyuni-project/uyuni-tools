package types

type PaygSshData struct {
	Description     string `mapstructure:"description"`
	Hostname        string `mapstructure:"hostname"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	BastionHostname string `mapstructure:"bastion_hostname"`
	BastionPort     int    `mapstructure:"bastion_port"`
	BastionUsername string `mapstructure:"bastion_username"`
}
