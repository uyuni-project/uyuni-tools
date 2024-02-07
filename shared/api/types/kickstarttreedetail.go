











package types

type KickstartTreeDetail struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    AbsPath string `mapstructure:"abs_path"`
    ChannelId int `mapstructure:"channel_id"`
    KernelOptions string `mapstructure:"kernel_options"`
    PostKernelOptions string `mapstructure:"post_kernel_options"`
	InstallType struct {
    Id int `mapstructure:"id"`
    Label string `mapstructure:"label"`
    Name string `mapstructure:"name"`
} `mapstructure:"install_type"`
} 