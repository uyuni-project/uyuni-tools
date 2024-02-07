package types

type RecurringAction struct {
	Id         int      `mapstructure:"id"`
	Name       string   `mapstructure:"name"`
	EntityId   int      `mapstructure:"entity_id"`
	EntityType string   `mapstructure:"entity_type"`
	CronExpr   string   `mapstructure:"cron_expr"`
	Created    string   `mapstructure:"created"`
	Creator    string   `mapstructure:"creator"`
	Test       bool     `mapstructure:"test"`
	States     []string `mapstructure:"states"`
	Active     bool     `mapstructure:"active"`
}
