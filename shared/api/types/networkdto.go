











package types

type NetworkDto struct {
    SystemId int `mapstructure:"systemId"`
    SystemName string `mapstructure:"systemName"`
    LastCheckin string `mapstructure:"last_checkin"`
} 