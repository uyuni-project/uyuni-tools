











package types

type ContentSource struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    SourceUrl string `mapstructure:"sourceUrl"`
    Type string `mapstructure:"type"`
    HasSignedMetadata bool `mapstructure:"hasSignedMetadata"`
	SslContentSources struct {
    SslCaDesc string `mapstructure:"sslCaDesc"`
    SslCertDesc string `mapstructure:"sslCertDesc"`
    SslKeyDesc string `mapstructure:"sslKeyDesc"`
} `mapstructure:"sslContentSources"`
} 