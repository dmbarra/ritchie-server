package starter

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"ritchie-server/server"
	"ritchie-server/server/config"
	"ritchie-server/server/http/cliversion"
	"ritchie-server/server/http/credential"
	"ritchie-server/server/http/formulas"
	"ritchie-server/server/http/health"
	"ritchie-server/server/http/hello"
	configHttp "ritchie-server/server/http/keycloak"
	"ritchie-server/server/http/login"
	"ritchie-server/server/http/oauth"
	"ritchie-server/server/http/repository"
	"ritchie-server/server/http/tree"
	"ritchie-server/server/http/usagelogger"
	"ritchie-server/server/http/user"
	"ritchie-server/server/keycloak"
	"ritchie-server/server/logger"
	"ritchie-server/server/middleware"
	"ritchie-server/server/ph"
	"ritchie-server/server/security"
	"ritchie-server/server/vault"
)

const fileSecurityConstraints string = "./server/resources/security-constraints.yml"

var fileConfig = os.Getenv(server.FileConfig)

type Configurator struct {
	conf            server.Config
	vaultManager    server.VaultManager
	keycloakManager server.KeycloakManager
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
	conf := config.NewConfiguration(loadConfigs(), sc)
	vm := vault.NewVaultManager(client)
	km := keycloak.NewKeycloakManager(conf)
	vm.Start(client)
	return Configurator{
		conf:            conf,
		vaultManager:    vm,
		keycloakManager: km,
	}, nil
}

func (c Configurator) LoadLoginHandler() server.DefaultHandler {
	return login.NewLoginHandler(c.keycloakManager)
}

func (c Configurator) LoadConfigHandler() server.DefaultHandler {
	return configHttp.NewKeycloakHandler(c.conf)
}

func (c Configurator) LoadUserHandler() server.DefaultHandler {
	return user.NewUserHandler(c.keycloakManager, c.vaultManager)
}

func (c Configurator) LoadCredentialConfigHandler() server.DefaultHandler {
	return credential.NewConfigHandler(c.conf)
}

func (c Configurator) LoadConfigHealth() server.DefaultHandler {
	return health.NewConfigHealth(c.conf)
}

func (c Configurator) LoadOauthHandler() server.DefaultHandler {
	return oauth.NewConfigHandler(c.conf)
}

func (c Configurator) LoadUsageLoggerHandler() server.DefaultHandler {
	return usagelogger.NewUsageLoggerHandler(c.conf)
}

func (c Configurator) LoadCliVersionHandler() server.DefaultHandler {
	return cliversion.NewConfigHandler(c.conf)
}

func (c Configurator) LoadRepositoryHandler() server.DefaultHandler {
	return repository.NewConfigHandler(c.conf)
}

func (c Configurator) LoadTreeHandler() server.DefaultHandler {
	sa := security.NewAuthorization(c.conf)
	ph := ph.NewProviderHandler(sa)
	return tree.NewConfigHandler(c.conf, sa, ph)
}

func (c Configurator) LoadFormulasHandler() server.DefaultHandler {
	sa := security.NewAuthorization(c.conf)
	ph := ph.NewProviderHandler(sa)
	return formulas.NewConfigHandler(c.conf, sa, ph)
}

func (c Configurator) LoadMiddlewareHandler() server.MiddlewareHandler {
	sa := security.NewAuthorization(c.conf)
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
