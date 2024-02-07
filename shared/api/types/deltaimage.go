











package types

type DeltaImage struct {
    SourceId int `mapstructure:"source_id"`
    TargetId int `mapstructure:"target_id"`
    File string `mapstructure:"file"`
    Pillar struct `mapstructure:"pillar"`
} 