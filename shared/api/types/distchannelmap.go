package types

type DistChannelMap struct {
	Os           string `mapstructure:"os"`
	Release      string `mapstructure:"release"`
	ArchName     string `mapstructure:"arch_name"`
	ChannelLabel string `mapstructure:"channel_label"`
	OrgSpecific  string `mapstructure:"org_specific"`
}
