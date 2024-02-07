package types

type KickstartDto struct {
	Label        string `mapstructure:"label"`
	TreeLabel    string `mapstructure:"tree_label"`
	Name         string `mapstructure:"name"`
	AdvancedMode bool   `mapstructure:"advanced_mode"`
	OrgDefault   bool   `mapstructure:"org_default"`
	Active       bool   `mapstructure:"active"`
	UpdateType   string `mapstructure:"update_type"`
}
