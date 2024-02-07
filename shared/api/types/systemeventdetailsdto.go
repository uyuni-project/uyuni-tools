package types

type SystemEventDetailsDto struct {
	Id          int    `mapstructure:"id"`
	HistoryType string `mapstructure:"history_type"`
	Status      string `mapstructure:"status"`
	Summary     string `mapstructure:"summary"`

	Created   string `mapstructure:"created"`
	PickedUp  string `mapstructure:"picked_up"`
	Completed string `mapstructure:"completed"`

	EarliestAction string `mapstructure:"earliest_action"`
	ResultMsg      string `mapstructure:"result_msg"`
	ResultCode     int    `mapstructure:"result_code"`
	AdditionalInfo struct {
		Detail string `mapstructure:"detail"`
		Result string `mapstructure:"result"`
	} `mapstructure:"additional_info"`
}
