package types

type Server struct {
	Id              int    `mapstructure:"id"`
	ProfileName     string `mapstructure:"profile_name"`
	MachineId       string `mapstructure:"machine_id"`
	Payg            bool   `mapstructure:"payg"`
	MinionId        string `mapstructure:"minion_id"`
	BaseEntitlement string `mapstructure:"base_entitlement"`

	AddonEntitlements []string `mapstructure:"addon_entitlements"`
	AutoUpdate        bool     `mapstructure:"auto_update"`
	Release           string   `mapstructure:"release"`
	Address1          string   `mapstructure:"address1"`
	Address2          string   `mapstructure:"address2"`
	City              string   `mapstructure:"city"`
	State             string   `mapstructure:"state"`
	Country           string   `mapstructure:"country"`
	Building          string   `mapstructure:"building"`
	Room              string   `mapstructure:"room"`
	Rack              string   `mapstructure:"rack"`
	Description       string   `mapstructure:"description"`
	Hostname          string   `mapstructure:"hostname"`
	LastBoot          string   `mapstructure:"last_boot"`
	OsaStatus         string   `mapstructure:"osa_status"`
	LockStatus        bool     `mapstructure:"lock_status"`
	Virtualization    string   `mapstructure:"virtualization"`
	ContactMethod     string   `mapstructure:"contact_method"`
}
