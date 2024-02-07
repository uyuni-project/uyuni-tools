package types

type CVEAuditServer struct {
	SystemId    int                 `mapstructure:"system_id"`
	PatchStatus string              `mapstructure:"patch_status"`
	String      []channel_labels    `mapstructure:"string"`
	String      []errata_advisories `mapstructure:"string"`
}
