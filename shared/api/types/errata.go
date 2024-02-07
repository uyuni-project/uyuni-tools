package types

type Errata struct {
	Id               int    `mapstructure:"id"`
	Date             string `mapstructure:"date"`
	AdvisoryType     string `mapstructure:"advisory_type"`
	AdvisoryStatus   string `mapstructure:"advisory_status"`
	AdvisoryName     string `mapstructure:"advisory_name"`
	AdvisorySynopsis string `mapstructure:"advisory_synopsis"`
}
