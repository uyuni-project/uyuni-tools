











package types

type TrustedOrgDto struct {
    OrgId int `mapstructure:"org_id"`
    OrgName string `mapstructure:"org_name"`
    SharedChannels int `mapstructure:"shared_channels"`
} 