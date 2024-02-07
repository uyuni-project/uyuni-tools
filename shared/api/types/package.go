











package types

type Package struct {
    Name string `mapstructure:"name"`
    Version string `mapstructure:"version"`
    Release string `mapstructure:"release"`
    Epoch string `mapstructure:"epoch"`
    Id int `mapstructure:"id"`
    ArchLabel string `mapstructure:"arch_label"`
    LastModified string `mapstructure:"last_modified"`
    Path string `mapstructure:"path"`
    PartOfRetractedPatch bool `mapstructure:"part_of_retracted_patch"`
    Provider string `mapstructure:"provider"`
} 