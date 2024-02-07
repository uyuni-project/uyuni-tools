package types

type ShortSystemInfo struct {
	Id          int    `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	LastCheckin string `mapstructure:"last_checkin"`
	Created     string `mapstructure:"created"`
	LastBoot    string `mapstructure:"last_boot"`
}
