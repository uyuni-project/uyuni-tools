package types

type XccdfTestResultDto struct {
	Xid       int    `mapstructure:"xid"`
	Profile   string `mapstructure:"profile"`
	Path      string `mapstructure:"path"`
	Ovalfiles string `mapstructure:"ovalfiles"`
	Completed string `mapstructure:"completed"`
}
