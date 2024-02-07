package types

type ConfigFileNameDto struct {
	Type         string `mapstructure:"type"`
	Path         string `mapstructure:"path"`
	ChannelLabel string `mapstructure:"channel_label"`
	ChannelType  struct {
		Id       int    `mapstructure:"id"`
		Label    string `mapstructure:"label"`
		Name     string `mapstructure:"name"`
		Priority int    `mapstructure:"priority"`
	} `mapstructure:"channel_type"`
	LastModified string `mapstructure:"last_modified"`
}
