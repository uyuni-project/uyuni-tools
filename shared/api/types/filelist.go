











package types

type FileList struct {
    Name string `mapstructure:"name"`
    FileNames []string `mapstructure:"file_names"`
} 