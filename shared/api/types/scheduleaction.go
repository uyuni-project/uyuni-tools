package types

type ScheduleAction struct {
	Id                int    `mapstructure:"id"`
	Name              string `mapstructure:"name"`
	Type              string `mapstructure:"type"`
	Scheduler         string `mapstructure:"scheduler"`
	Earliest          string `mapstructure:"earliest"`
	Prerequisite      int    `mapstructure:"prerequisite"`
	CompletedSystems  int    `mapstructure:"completedSystems"`
	FailedSystems     int    `mapstructure:"failedSystems"`
	InProgressSystems int    `mapstructure:"inProgressSystems"`
}
