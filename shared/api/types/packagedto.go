package types

type PackageDto struct {
	Name             string `mapstructure:"name"`
	Version          string `mapstructure:"version"`
	Release          string `mapstructure:"release"`
	Epoch            string `mapstructure:"epoch"`
	Checksum         string `mapstructure:"checksum"`
	ChecksumType     string `mapstructure:"checksum_type"`
	Id               int    `mapstructure:"id"`
	ArchLabel        string `mapstructure:"arch_label"`
	LastModifiedDate string `mapstructure:"last_modified_date"`
	LastModified     string `mapstructure:"last_modified"`
}
