











package types

type VirtualSystemOverview struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    GuestName string `mapstructure:"guest_name"`
    LastCheckin string `mapstructure:"last_checkin"`
    Uuid string `mapstructure:"uuid"`
} 