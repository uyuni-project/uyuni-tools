











package types

type FormulaData struct {
    SystemId int `mapstructure:"system_id"`
    MinionId string `mapstructure:"minion_id"`
    FormulaValues struct `mapstructure:"formula_values"`
} 