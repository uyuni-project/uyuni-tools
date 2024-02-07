











package types

type ConfigFileDto struct {
    Type string `mapstructure:"type"`
    Path string `mapstructure:"path"`
    LastModified string `mapstructure:"last_modified"`
} 