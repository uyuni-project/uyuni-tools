package types

type IssMasterOrg struct {
	MasterOrgId   int    `mapstructure:"masterOrgId"`
	MasterOrgName string `mapstructure:"masterOrgName"`
	LocalOrgId    int    `mapstructure:"localOrgId"`
}
