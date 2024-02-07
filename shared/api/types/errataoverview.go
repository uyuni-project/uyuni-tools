package types

type ErrataOverview struct {
	Id               int    `mapstructure:"id"`
	IssueDate        string `mapstructure:"issue_date"`
	Date             string `mapstructure:"date"`
	UpdateDate       string `mapstructure:"update_date"`
	AdvisorySynopsis string `mapstructure:"advisory_synopsis"`
	AdvisoryType     string `mapstructure:"advisory_type"`
	AdvisoryStatus   string `mapstructure:"advisory_status"`
	AdvisoryName     string `mapstructure:"advisory_name"`
}
