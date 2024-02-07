











package types

type PackageState struct {
    Id int `mapstructure:"id"`
    Name string `mapstructure:"name"`
    StateRevisionId int `mapstructure:"state_revision_id"`
    PackageStateTypeId string `mapstructure:"package_state_type_id"`
    VersionConstraintId string `mapstructure:"version_constraint_id"`
} 