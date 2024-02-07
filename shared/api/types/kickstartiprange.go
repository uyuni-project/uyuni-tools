package types

type KickstartIpRange struct {
	KsLabel string `mapstructure:"ksLabel"`
	Max     string `mapstructure:"max"`
	Min     string `mapstructure:"min"`
}
