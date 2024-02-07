











package types

type ServerSnapshot struct {
    Id int `mapstructure:"id"`
    Reason string `mapstructure:"reason"`
    Created string `mapstructure:"created"`
    Channels []string `mapstructure:"channels"`
    Groups []string `mapstructure:"groups"`
    Entitlements []string `mapstructure:"entitlements"`
    ConfigChannels []string `mapstructure:"config_channels"`
    Tags []string `mapstructure:"tags"`
    InvalidReason string `mapstructure:"Invalid_reason"`
} 