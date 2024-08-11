package config

type ServiceConfig struct {
	DatabaseURL      string `yaml:"database_url"`
	RegistryURL      string `yaml:"registry_url"`
	ServiceName      string `yaml:"service_name"`
	ListeningAddress string `yaml:"listening_address"`
}
