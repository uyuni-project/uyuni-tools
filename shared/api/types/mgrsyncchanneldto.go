











package types

type MgrSyncChannelDto struct {
    Arch string `mapstructure:"arch"`
    Description string `mapstructure:"description"`
    Family string `mapstructure:"family"`
    IsSigned bool `mapstructure:"is_signed"`
    Label string `mapstructure:"label"`
    Name string `mapstructure:"name"`
    Optional bool `mapstructure:"optional"`
    Parent string `mapstructure:"parent"`
    ProductName string `mapstructure:"product_name"`
    ProductVersion string `mapstructure:"product_version"`
    SourceUrl string `mapstructure:"source_url"`
    Status string `mapstructure:"status"`
    Summary string `mapstructure:"summary"`
    UpdateTag string `mapstructure:"update_tag"`
    InstallerUpdates bool `mapstructure:"installer_updates"`
} 