package main

type MapUserConfig struct {
	UserArn string `yaml:"userarn"`
	Username string `yaml:"username"`
	Groups []string `yaml:"groups"`
}
