package types

type KickstartOptionValue struct {
	Name    string `mapstructure:"name"`
	Value   string `mapstructure:"value"`
	Enabled bool   `mapstructure:"enabled"`
}
