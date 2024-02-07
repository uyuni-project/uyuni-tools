











package types

type SUSEInstalledProduct struct {
    Name string `mapstructure:"name"`
    IsBaseProduct bool `mapstructure:"isBaseProduct"`
    Version string `mapstructure:"version"`
    Arch string `mapstructure:"arch"`
    Release string `mapstructure:"release"`
    FriendlyName string `mapstructure:"friendlyName"`
} 