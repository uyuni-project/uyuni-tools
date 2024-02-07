package types

type CVEAuditImage struct {
	ImageId     int                 `mapstructure:"image_id"`
	PatchStatus string              `mapstructure:"patch_status"`
	String      []channel_labels    `mapstructure:"string"`
	String      []errata_advisories `mapstructure:"string"`
}
