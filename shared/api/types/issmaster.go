











package types

type IssMaster struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    CaCert string `mapstructure:"caCert"`
    IsCurrentMaster bool `mapstructure:"isCurrentMaster"`
} 