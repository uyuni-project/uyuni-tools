package types

type Formula struct {
	Name         string `mapstructure:"name"`
	Description  string `mapstructure:"description"`
	FormulaGroup string `mapstructure:"formula_group"`
}
