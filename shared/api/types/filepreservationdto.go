











package types

type FilePreservationDto struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Created string `mapstructure:"created"`
    LastModified string `mapstructure:"last_modified"`
} 