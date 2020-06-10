package config

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"ritchie-server/server"
)

type Configuration struct {
	Configs             map[string]*server.ConfigFile
	SecurityConstraints server.SecurityConstraints
}

func NewConfiguration(c map[string]*server.ConfigFile, s server.SecurityConstraints) server.Config {
	return Configuration{Configs: c, SecurityConstraints: s}
}

func (c Configuration) ReadHealthConfigs() map[string]server.HealthEndpoints {

	m := make(map[string]server.HealthEndpoints)

	for orgName := range c.Configs {
		if orgName != "default" {
			vaultConfig := api.DefaultConfig()
			_ = vaultConfig.ReadEnvironment()

			url := fmt.Sprint(vaultConfig.Address, "/sys/health")
			m[orgName] = server.HealthEndpoints{
				VaultURL:    url,
			}
		}
	}
	return m
}

func (c Configuration) ReadCredentialConfigs(org string) (map[string][]server.CredentialConfig, error) {
	config, err := c.getConfig(org)
	if err != nil  {
		return nil, err
	}
	return config.CredentialConfig, nil
}

func (c Configuration) ReadCliVersionConfigs(org string) (server.CliVersionConfig, error) {
	config, err := c.getConfig(org)
	if err != nil {
		return server.CliVersionConfig{}, err
	}

	return config.CliVersionConfig, nil
}

func (c Configuration) ReadRepositoryConfig(org string) ([]server.Repository, error) {
	config, err := c.getConfig(org)
	if err != nil {
		return nil, err
	}

	return config.RepositoryConfig, nil
}

func (c Configuration) ReadSecurityConstraints() server.SecurityConstraints {
	return c.SecurityConstraints
}

func (c Configuration) getConfig(org string) (*server.ConfigFile, error) {
	config := c.Configs[org]
	if config == nil {
		return nil, fmt.Errorf("config not found for organization: %s", org)
	}

	return config, nil
}
