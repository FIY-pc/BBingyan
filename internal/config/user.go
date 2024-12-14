package config

type UserConfig struct {
	InitAdmin *InitAdminConfig `yaml:"initAdmin"`
}

type InitAdminConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}
