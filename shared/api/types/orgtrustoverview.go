











package types

type OrgTrustOverview struct {
    OrgId int `mapstructure:"orgId"`
    OrgName string `mapstructure:"orgName"`
    TrustEnabled bool `mapstructure:"trustEnabled"`
} 