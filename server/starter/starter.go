package starter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"ritchie-server/server"
	"ritchie-server/server/config"
	"ritchie-server/server/http/cliversion"
	"ritchie-server/server/http/credential"
	"ritchie-server/server/http/formulas"
	"ritchie-server/server/http/health"
	"ritchie-server/server/http/hello"
	"ritchie-server/server/http/login"
	"ritchie-server/server/http/repository"
	"ritchie-server/server/http/tree"
	"ritchie-server/server/http/ul"
	"ritchie-server/server/logger"
	"ritchie-server/server/middleware"
	"ritchie-server/server/fph"
	"ritchie-server/server/security"
	"ritchie-server/server/sp/keycloak"
	"ritchie-server/server/sp/ldap"
	"ritchie-server/server/vault"
)

const fileSecurityConstraints string = "./server/resources/security-constraints.yml"

var fileConfig = os.Getenv(server.FileConfig)

type Configurator struct {
	conf              server.Config
	vaultManager      server.VaultManager
	securityProviders server.SecurityProviders
}

func NewConfiguration() (server.Configurator, error) {

	logger.LoadLogDefinition()
	yml, _ := ioutil.ReadFile(fileSecurityConstraints)
	sc := server.SecurityConstraints{}
	err := yaml.Unmarshal(yml, &sc)
	if err != nil {
		return Configurator{}, fmt.Errorf("load security constraints error: %v", err)
	}
	client, err := vault.NewConfig().Start()
	if err != nil {
		return Configurator{}, fmt.Errorf("could not connect to vault, error: %v", err)
	}
	configs := loadConfigs()
	conf := config.NewConfiguration(configs, sc)
	vm := vault.NewVaultManager(client)
	vm.Start(client)
	sp := loadSecurityProviders(configs)

	return Configurator{
		conf:              conf,
		vaultManager:      vm,
		securityProviders: sp,
	}, nil
}

func (c Configurator) LoadLoginHandler() server.DefaultHandler {
	return login.NewLoginHandler(c.securityProviders, c.vaultManager)
}

func (c Configurator) LoadCredentialConfigHandler() server.DefaultHandler {
	return credential.NewConfigHandler(c.conf)
}

func (c Configurator) LoadConfigHealth() server.DefaultHandler {
	return health.NewConfigHealth(c.conf)
}

func (c Configurator) LoadUsageLoggerHandler() server.DefaultHandler {
	return ul.NewUsageLoggerHandler()
}

func (c Configurator) LoadCliVersionHandler() server.DefaultHandler {
	return cliversion.NewConfigHandler(c.conf)
}

func (c Configurator) LoadRepositoryHandler() server.DefaultHandler {
	return repository.NewConfigHandler(c.conf)
}

func (c Configurator) LoadTreeHandler() server.DefaultHandler {
	sa := security.NewAuthorization(c.conf, c.vaultManager)
	ph := fph.NewProviderHandler(sa)
	return tree.NewConfigHandler(c.conf, sa, ph)
}

func (c Configurator) LoadFormulasHandler() server.DefaultHandler {
	sa := security.NewAuthorization(c.conf, c.vaultManager)
	ph := fph.NewProviderHandler(sa)
	return formulas.NewConfigHandler(c.conf, sa, ph)
}

func (c Configurator) LoadMiddlewareHandler() server.MiddlewareHandler {
	sa := security.NewAuthorization(c.conf, c.vaultManager)
	return middleware.NewMiddlewareHandler(sa)
}

func (c Configurator) LoadCredentialHandler() server.CredentialHandler {
	return credential.NewCredentialHandler(c.vaultManager, c.conf)
}

func (c Configurator) LoadHelloHandler() server.DefaultHandler {
	return hello.NewHelloHandler()
}

func loadConfigs() map[string]*server.ConfigFile {
	config, err := readFileConfig()
	if err != nil {
		log.Fatal("Load configs error ", err)
		return nil
	}

	return config
}

func readFileConfig() (map[string]*server.ConfigFile, error) {
	file, err := ioutil.ReadFile(fileConfig)
	if err != nil {
		log.Error("Read file config error ", err)
		return nil, err
	}

	var configFile map[string]*server.ConfigFile
	_ = json.Unmarshal(file, &configFile)
	return configFile, nil
}

func loadSecurityProviders(config map[string]*server.ConfigFile) server.SecurityProviders {
	sm := make(map[string]server.SecurityManager)
	for org, config := range config {
		switch config.SPConfig["type"] {
		case "keycloak":
			sm[org] = keycloak.NewKeycloakProvider(config.SPConfig)
		case "ldap":
			sm[org] = ldap.NewLdapProvider(config.SPConfig)

		}
	}
	return server.SecurityProviders{Providers: sm}
}
