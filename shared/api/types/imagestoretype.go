package types

type ImageStoreType struct {
	Id    int    `mapstructure:"id"`
	Label string `mapstructure:"label"`
	Name  string `mapstructure:"name"`
}
