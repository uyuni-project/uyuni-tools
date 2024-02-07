











package types

type MaintenanceSchedule struct {
    Id int `mapstructure:"id"`
    OrgId int `mapstructure:"orgId"`
    Name string `mapstructure:"name"`
    Type string `mapstructure:"type"`
	Calendar struct {
    Id int `mapstructure:"id"`
    OrgId int `mapstructure:"orgId"`
    Label string `mapstructure:"label"`
    Url string `mapstructure:"url"`
    Ical string `mapstructure:"ical"`
} `mapstructure:"calendar"`
} 