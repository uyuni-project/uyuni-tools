











package types

type RhnTimeZone struct {
    TimeZoneId int `mapstructure:"time_zone_id"`
    OlsonName string `mapstructure:"olson_name"`
} 