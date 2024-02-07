package types

type ContentProjectSource struct {
	ContentProjectLabel string `mapstructure:"contentProjectLabel"`
	Type                string `mapstructure:"type"`
	State               string `mapstructure:"state"`
	ChannelLabel        string `mapstructure:"channelLabel"`
}
