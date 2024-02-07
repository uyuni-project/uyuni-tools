











package types

type KickstartTree struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    BasePath string `mapstructure:"base_path"`
    ChannelId int `mapstructure:"channel_id"`
} 