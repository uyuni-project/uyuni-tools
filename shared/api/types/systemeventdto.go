package types

type SystemEventDto struct {
	Id          int    `mapstructure:"id"`
	HistoryType string `mapstructure:"history_type"`
	Status      string `mapstructure:"status"`
	Summary     string `mapstructure:"summary"`
	Completed   string `mapstructure:"completed"`
}
