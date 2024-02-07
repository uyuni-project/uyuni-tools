











package types

type ImageStore struct {
    Label string `mapstructure:"label"`
    Uri string `mapstructure:"uri"`
    Storetype string `mapstructure:"storetype"`
    HasCredentials bool `mapstructure:"hasCredentials"`
    Username string `mapstructure:"username"`
} 