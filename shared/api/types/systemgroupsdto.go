package types

type SystemGroupsDTO struct {
	Id           int `mapstructure:"id"`
	SystemGroups struct {
		Id   int    `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"system_groups"`
}
