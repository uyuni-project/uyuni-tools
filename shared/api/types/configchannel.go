package types

type ConfigChannel struct {
	Id                int    `mapstructure:"id"`
	OrgId             int    `mapstructure:"orgId"`
	Label             string `mapstructure:"label"`
	Name              string `mapstructure:"name"`
	Description       string `mapstructure:"description"`
	ConfigChannelType struct {
		Id       int    `mapstructure:"id"`
		Label    string `mapstructure:"label"`
		Name     string `mapstructure:"name"`
		Priority int    `mapstructure:"priority"`
	} `mapstructure:"configChannelType"`
}
