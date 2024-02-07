package types

type CustomDataKey struct {
	Id           int    `mapstructure:"id"`
	Label        string `mapstructure:"label"`
	Description  string `mapstructure:"description"`
	SystemCount  int    `mapstructure:"system_count"`
	LastModified string `mapstructure:"last_modified"`
}
