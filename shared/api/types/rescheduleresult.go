











package types

type RescheduleResult struct {
    Strategy string `mapstructure:"strategy"`
    ForScheduleName string `mapstructure:"for_schedule_name"`
    Status bool `mapstructure:"status"`
    Message string `mapstructure:"message"`
	Actions struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    Type string `mapstructure:"type"`
    Scheduler string `mapstructure:"scheduler"`
    Earliest string `mapstructure:"earliest"`
    Prerequisite int `mapstructure:"prerequisite"`
    AffectedSystemIds []int `mapstructure:"affected_system_ids"`
    Details string `mapstructure:"details"`
} `mapstructure:"actions"`
} 