package types

type ProfileOverviewDto struct {
	Id      int    `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Channel string `mapstructure:"channel"`
}
