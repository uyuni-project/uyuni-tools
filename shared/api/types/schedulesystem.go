











package types

type ScheduleSystem struct {
    ServerId int `mapstructure:"server_id"`
    ServerName string `mapstructure:"server_name"`
    BaseChannel string `mapstructure:"base_channel"`
    Timestamp string `mapstructure:"timestamp"`
    Message string `mapstructure:"message"`
} 