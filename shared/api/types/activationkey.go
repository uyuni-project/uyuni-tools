











package types

type ActivationKey struct {
    Key string `mapstructure:"key"`
    Description string `mapstructure:"description"`
    UsageLimit int `mapstructure:"usage_limit"`
    BaseChannelLabel string `mapstructure:"base_channel_label"`
    ChildChannelLabels []string `mapstructure:"child_channel_labels"`
    Entitlements []string `mapstructure:"entitlements"`
    ServerGroupIds []string `mapstructure:"server_group_ids"`
    PackageNames []string `mapstructure:"package_names"`
	Packages struct {
    Name string `mapstructure:"name"`
    Arch string `mapstructure:"arch"`
} `mapstructure:"packages"`
    UniversalDefault bool `mapstructure:"universal_default"`
    Disabled bool `mapstructure:"disabled"`
    ContactMethod string `mapstructure:"contact_method"`
} 