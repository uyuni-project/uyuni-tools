package types

type PackageOverview struct {
	Id          int    `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	Summary     string `mapstructure:"summary"`
	Description string `mapstructure:"description"`
	Version     string `mapstructure:"version"`
	Release     string `mapstructure:"release"`
	Arch        string `mapstructure:"arch"`
	Epoch       string `mapstructure:"epoch"`
	Provider    string `mapstructure:"provider"`
}
