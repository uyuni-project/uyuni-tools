











package types

type OrgUserExtGroup struct {
    Name string `mapstructure:"name"`
    Groups []string `mapstructure:"groups"`
} 