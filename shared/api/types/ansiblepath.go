package types

type AnsiblePath struct {
	Id       int    `mapstructure:"id"`
	Type     string `mapstructure:"type"`
	ServerId int    `mapstructure:"server_id"`
	Path     string `mapstructure:"path"`
}
