











package types

type MgrSyncProductDto struct {
    FriendlyName string `mapstructure:"friendly_name"`
    Arch string `mapstructure:"arch"`
    Status string `mapstructure:"status"`
	Channels struct {
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
}     #array_end()
	Extensions struct {
    FriendlyName string `mapstructure:"friendly_name"`
    Arch string `mapstructure:"arch"`
    Status string `mapstructure:"status"`
	Channels struct {
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
}         #array_end()
}     #array_end()
    Recommended bool `mapstructure:"recommended"`
} 