











package types

type OrgDto struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    ActiveUsers int `mapstructure:"active_users"`
    Systems int `mapstructure:"systems"`
    Trusts int `mapstructure:"trusts"`
    SystemGroups int `mapstructure:"system_groups"`
    ActivationKeys int `mapstructure:"activation_keys"`
    KickstartProfiles int `mapstructure:"kickstart_profiles"`
    ConfigurationChannels int `mapstructure:"configuration_channels"`
    StagingContentEnabled bool `mapstructure:"staging_content_enabled"`
} 