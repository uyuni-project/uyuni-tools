











package types

type ImageInfo struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Version string `mapstructure:"version"`
    Revision int `mapstructure:"revision"`
    Arch string `mapstructure:"arch"`
    External bool `mapstructure:"external"`
    StoreLabel string `mapstructure:"storeLabel"`
    Checksum string `mapstructure:"checksum"`
    Obsolete string `mapstructure:"obsolete"`
} 