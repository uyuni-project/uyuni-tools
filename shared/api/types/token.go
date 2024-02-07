











package types

type Token struct {
    Description string `mapstructure:"description"`
    UsageLimit int `mapstructure:"usage_limit"`
    BaseChannelLabel string `mapstructure:"base_channel_label"`
    ChildChannelLabels []string `mapstructure:"child_channel_labels"`
    Entitlements []string `mapstructure:"entitlements"`
    ServerGroupIds []string `mapstructure:"server_group_ids"`
    PackageNames []string `mapstructure:"package_names"`
	Packages struct {
    String name `mapstructure:"string"`
    String arch `mapstructure:"string"`
} `mapstructure:"packages"`
    UniversalDefault bool `mapstructure:"universal_default"`
    Disabled bool `mapstructure:"disabled"`
} 