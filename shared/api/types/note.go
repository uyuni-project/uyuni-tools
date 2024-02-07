package types

type Note struct {
	Id       int    `mapstructure:"id"`
	Subject  string `mapstructure:"subject"`
	Note     string `mapstructure:"note"`
	SystemId int    `mapstructure:"system_id"`
	Creator  string `mapstructure:"creator"`
	Updated  date   `mapstructure:"updated"`
}
