











package types

type MirrorCredentialsDto struct {
    Id int `mapstructure:"id"`
    User string `mapstructure:"user"`
    IsPrimary bool `mapstructure:"isPrimary"`
} 