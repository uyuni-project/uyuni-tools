package types

type ImageProfile struct {
	Label         string `mapstructure:"label"`
	ImageType     string `mapstructure:"imageType"`
	ImageStore    string `mapstructure:"imageStore"`
	ActivationKey string `mapstructure:"activationKey"`
	Path          string `mapstructure:"path"`
}
