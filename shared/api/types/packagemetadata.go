











package types

type PackageMetadata struct {
    PackageNameId int `mapstructure:"package_name_id"`
    PackageName string `mapstructure:"package_name"`
    PackageEpoch string `mapstructure:"package_epoch"`
    PackageVersion string `mapstructure:"package_version"`
    PackageRelease string `mapstructure:"package_release"`
    PackageArch string `mapstructure:"package_arch"`
    ThisSystem string `mapstructure:"this_system"`
    OtherSystem string `mapstructure:"other_system"`
    Comparison int `mapstructure:"comparison"`
} 