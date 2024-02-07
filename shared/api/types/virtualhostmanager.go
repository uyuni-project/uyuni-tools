











package types

type VirtualHostManager struct {
    Label string `mapstructure:"label"`
    OrgId int `mapstructure:"org_id"`
    GathererModule string `mapstructure:"gatherer_module"`
    Configs struct `mapstructure:"configs"`
} 