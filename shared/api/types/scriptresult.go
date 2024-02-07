











package types

type ScriptResult struct {
    ServerId int `mapstructure:"serverId"`
    StartDate string `mapstructure:"startDate"`
    StopDate string `mapstructure:"stopDate"`
    ReturnCode int `mapstructure:"returnCode"`
    Output string `mapstructure:"output"`
    OutputEnc64 bool `mapstructure:"output_enc64"`
} 