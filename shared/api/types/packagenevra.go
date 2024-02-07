











package types

type PackageNevra struct {
    Name string `mapstructure:"name"`
    Epoch string `mapstructure:"epoch"`
    Version string `mapstructure:"version"`
    Release string `mapstructure:"release"`
    Arch string `mapstructure:"arch"`
} 