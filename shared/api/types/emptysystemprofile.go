











package types

type EmptySystemProfile struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Created string `mapstructure:"created"`
    HwAddress []string `mapstructure:"hw_address"`
} 