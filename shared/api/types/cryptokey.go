











package types

type CryptoKey struct {
    Description string `mapstructure:"description"`
    Type string `mapstructure:"type"`
    Content string `mapstructure:"content"`
} 