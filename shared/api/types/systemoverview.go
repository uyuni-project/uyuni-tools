package types

type SystemOverview struct {
	Id               int    `mapstructure:"id"`
	Name             string `mapstructure:"name"`
	LastCheckin      string `mapstructure:"last_checkin"`
	Created          string `mapstructure:"created"`
	LastBoot         string `mapstructure:"last_boot"`
	ExtraPkgCount    int    `mapstructure:"extra_pkg_count"`
	OutdatedPkgCount int    `mapstructure:"outdated_pkg_count"`
}
