package types

type AnsiblePlaybook struct {
	Fullpath        string `mapstructure:"fullpath"`
	CustomInventory string `mapstructure:"custom_inventory"`
}
