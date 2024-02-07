











package types

type ContentEnvironment struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    Name string `mapstructure:"name"`
    Description string `mapstructure:"description"`
    Version int `mapstructure:"version"`
    Status string `mapstructure:"status"`
    LastBuildDate string `mapstructure:"lastBuildDate"`
    ContentProjectLabel string `mapstructure:"contentProjectLabel"`
    PreviousEnvironmentLabel string `mapstructure:"previousEnvironmentLabel"`
    NextEnvironmentLabel string `mapstructure:"nextEnvironmentLabel"`
} 