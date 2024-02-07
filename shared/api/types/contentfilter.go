











package types

type ContentFilter struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    OrgId int `mapstructure:"orgId"`
    EntityType string `mapstructure:"entityType"`
    Rule string `mapstructure:"rule"`
	Criteria struct {
    Matcher string `mapstructure:"matcher"`
    Field string `mapstructure:"field"`
    Value string `mapstructure:"value"`
} `mapstructure:"criteria"`
} 