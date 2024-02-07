











package types

type ManagedServerGroup struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Description string `mapstructure:"description"`
    OrgId int `mapstructure:"org_id"`
    SystemCount int `mapstructure:"system_count"`
} 