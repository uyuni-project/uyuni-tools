











package types

type ContentProject struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    Name string `mapstructure:"name"`
    Description string `mapstructure:"description"`
    LastBuildDate string `mapstructure:"lastBuildDate"`
    OrgId int `mapstructure:"orgId"`
    FirstEnvironment string `mapstructure:"firstEnvironment"`
} 