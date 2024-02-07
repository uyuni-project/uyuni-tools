











package types

type ConfigRevision struct {
    Type string `mapstructure:"type"`
    Path string `mapstructure:"path"`
    TargetPath string `mapstructure:"target_path"`
    Channel string `mapstructure:"channel"`
    Contents string `mapstructure:"contents"`
    ContentsEnc64 bool `mapstructure:"contents_enc64"`
    Revision int `mapstructure:"revision"`
    Creation string `mapstructure:"creation"`
    Modified string `mapstructure:"modified"`
    Owner string `mapstructure:"owner"`
    Group string `mapstructure:"group"`
    Permissions int `mapstructure:"permissions"`
    PermissionsMode string `mapstructure:"permissions_mode"`
    SelinuxCtx string `mapstructure:"selinux_ctx"`
    Binary bool `mapstructure:"binary"`
    Sha256 string `mapstructure:"sha256"`
    MacroStartDelimiter string `mapstructure:"macro-start-delimiter"`
    MacroEndDelimiter string `mapstructure:"macro-end-delimiter"`
} 