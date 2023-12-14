package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Jira   Jira   `yaml:"jira"`
	Github Github `yaml:"github"`
}

type Jira struct {
	BaseURL  string `yaml:"base_url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Github struct {
	OAuthToken string `yaml:"oauth_token,omitempty"`
}

func Load(path string) (cfg *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&cfg)
	return
}
