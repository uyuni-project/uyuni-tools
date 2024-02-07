











package types

type IssSlave struct {
    Id int `mapstructure:"id"`
    Slave string `mapstructure:"slave"`
    Enabled bool `mapstructure:"enabled"`
    AllowAllOrgs bool `mapstructure:"allowAllOrgs"`
} 