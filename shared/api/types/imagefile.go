











package types

type ImageFile struct {
    File string `mapstructure:"file"`
    Type string `mapstructure:"type"`
    External bool `mapstructure:"external"`
    Url string `mapstructure:"url"`
} 