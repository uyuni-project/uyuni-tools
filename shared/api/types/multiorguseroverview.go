package types

type MultiOrgUserOverview struct {
	Login      string `mapstructure:"login"`
	LoginUc    string `mapstructure:"login_uc"`
	Name       string `mapstructure:"name"`
	Email      string `mapstructure:"email"`
	IsOrgAdmin bool   `mapstructure:"is_org_admin"`
}
