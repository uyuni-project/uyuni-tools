











package types

type ChannelTreeNode struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    Name string `mapstructure:"name"`
    ProviderName string `mapstructure:"provider_name"`
    Packages int `mapstructure:"packages"`
    Systems int `mapstructure:"systems"`
    ArchName string `mapstructure:"arch_name"`
} 