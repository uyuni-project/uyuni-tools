











package types

type EndpointInfo struct {
    SystemId int `mapstructure:"system_id"`
    EndpointName string `mapstructure:"endpoint_name"`
    ExporterName string `mapstructure:"exporter_name"`
    Module string `mapstructure:"module"`
    Path string `mapstructure:"path"`
    Port int `mapstructure:"port"`
    TlsEnabled bool `mapstructure:"tls_enabled"`
} 