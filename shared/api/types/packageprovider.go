











package types

type PackageProvider struct {
    Name string `mapstructure:"name"`
	Keys struct {
    Key string `mapstructure:"key"`
    Type string `mapstructure:"type"`
} `mapstructure:"keys"`
} 