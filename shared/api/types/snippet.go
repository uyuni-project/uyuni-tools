











package types

type Snippet struct {
    Name string `mapstructure:"name"`
    Contents string `mapstructure:"contents"`
    Fragment string `mapstructure:"fragment"`
    File string `mapstructure:"file"`
} 