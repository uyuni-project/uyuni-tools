











package types

type ContentSourceFilter struct {
    SortOrder int `mapstructure:"sortOrder"`
    Filter string `mapstructure:"filter"`
    Flag string `mapstructure:"flag"`
} 