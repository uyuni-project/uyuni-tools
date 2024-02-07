











package types

type HistoryEvent struct {
    Completed string `mapstructure:"completed"`
    Summary string `mapstructure:"summary"`
    Details string `mapstructure:"details"`
} 