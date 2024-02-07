











package types

type ImageOverview struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Type string `mapstructure:"type"`
    Version string `mapstructure:"version"`
    Revision int `mapstructure:"revision"`
    Arch string `mapstructure:"arch"`
    External bool `mapstructure:"external"`
    Checksum string `mapstructure:"checksum"`
    ProfileLabel string `mapstructure:"profileLabel"`
    StoreLabel string `mapstructure:"storeLabel"`
    BuildStatus string `mapstructure:"buildStatus"`
    InspectStatus string `mapstructure:"inspectStatus"`
    BuildServerId int `mapstructure:"buildServerId"`
    SecurityErrata int `mapstructure:"securityErrata"`
    BugErrata int `mapstructure:"bugErrata"`
    EnhancementErrata int `mapstructure:"enhancementErrata"`
    OutdatedPackages int `mapstructure:"outdatedPackages"`
    InstalledPackages int `mapstructure:"installedPackages"`
    Files struct `mapstructure:"files"`
    Obsolete bool `mapstructure:"obsolete"`
} 