











package types

type Channel struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Label string `mapstructure:"label"`
    ArchName string `mapstructure:"arch_name"`
    ArchLabel string `mapstructure:"arch_label"`
    Summary string `mapstructure:"summary"`
    Description string `mapstructure:"description"`
    ChecksumLabel string `mapstructure:"checksum_label"`
    LastModified string `mapstructure:"last_modified"`
    MaintainerName string `mapstructure:"maintainer_name"`
    MaintainerEmail string `mapstructure:"maintainer_email"`
    MaintainerPhone string `mapstructure:"maintainer_phone"`
    SupportPolicy string `mapstructure:"support_policy"`
    GpgKeyUrl string `mapstructure:"gpg_key_url"`
    GpgKeyId string `mapstructure:"gpg_key_id"`
    GpgKeyFp string `mapstructure:"gpg_key_fp"`
    YumrepoLastSync string `mapstructure:"yumrepo_last_sync"`
    EndOfLife string `mapstructure:"end_of_life"`
    ParentChannelLabel string `mapstructure:"parent_channel_label"`
    CloneOriginal string `mapstructure:"clone_original"`
	ContentSources struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    SourceUrl string `mapstructure:"sourceUrl"`
    Type string `mapstructure:"type"`
} `mapstructure:"contentSources"`
} 