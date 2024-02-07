package types

type ServerAction struct {
	FailedCount     int    `mapstructure:"failed_count"`
	Modified        string `mapstructure:"modified"`
	ModifiedDate    string `mapstructure:"modified_date"`
	Created         string `mapstructure:"created"`
	CreatedDate     string `mapstructure:"created_date"`
	ActionType      string `mapstructure:"action_type"`
	SuccessfulCount int    `mapstructure:"successful_count"`
	EarliestAction  string `mapstructure:"earliest_action"`
	Archived        int    `mapstructure:"archived"`
	SchedulerUser   string `mapstructure:"scheduler_user"`
	Prerequisite    string `mapstructure:"prerequisite"`
	Name            string `mapstructure:"name"`
	Id              int    `mapstructure:"id"`
	Version         string `mapstructure:"version"`
	CompletionTime  string `mapstructure:"completion_time"`
	CompletedDate   string `mapstructure:"completed_date"`
	PickupTime      string `mapstructure:"pickup_time"`
	PickupDate      string `mapstructure:"pickup_date"`
	ResultMsg       string `mapstructure:"result_msg"`
}
