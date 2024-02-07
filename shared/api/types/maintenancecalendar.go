











package types

type MaintenanceCalendar struct {
    Id int `mapstructure:"id"`
    OrgId int `mapstructure:"orgId"`
    Label string `mapstructure:"label"`
    Url string `mapstructure:"url"`
    Ical string `mapstructure:"ical"`
} 