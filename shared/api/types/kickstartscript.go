package types

type KickstartScript struct {
	Id                 int    `mapstructure:"id"`
	Name               string `mapstructure:"name"`
	Contents           string `mapstructure:"contents"`
	ScriptType         string `mapstructure:"script_type"`
	Interpreter        string `mapstructure:"interpreter"`
	Chroot             bool   `mapstructure:"chroot"`
	Erroronfail        bool   `mapstructure:"erroronfail"`
	Template           bool   `mapstructure:"template"`
	BeforeRegistration bool   `mapstructure:"beforeRegistration"`
}
