package server

import (
	"net/http"

	"github.com/hashicorp/vault/api"
)

const (
	OrganizationHeader = "x-org"
	ContextHeader      = "x-ctx"
	FileConfig         = "FILE_CONFIG"
	RepoNameHeader      = "x-repo-name"
	AuthorizationHeader = "x-authorization"
)

type (
	Org string
	Ctx string

	Repository struct {
		Name           string   `json:"name"`
		Priority       int      `json:"priority"`
		TreePath       string   `json:"treePath"`
		ServerUrl      string   `json:"serverUrl,omitempty"`
		ReplaceRepoUrl string   `json:"replaceRepoUrl,omitempty"`
		Username       string   `json:"username,omitempty"`
		Password       string   `json:"password,omitempty"`
		Provider       Provider `json:"provider,omitempty"`
	}
	Provider struct {
		Type   string `json:"type"`
		Bucket string `json:"bucket,omitempty"`
		Region string `json:"region,omitempty"`
		Remote string `json:"remote,omitempty"`
	}

	Tree struct {
		Commands []Command `json:"commands"`
		Version  string    `json:"version"`
	}

	Command struct {
		Usage   string   `json:"usage"`
		Help    string   `json:"help"`
		Formula *formula `json:"formula,omitempty"`
		Parent  string   `json:"parent"`
		Roles   []string `json:"roles,omitempty"`
	}

	formula struct {
		Path       string `json:"path"`
		Bin        string `json:"bin,omitempty"`
		BinWindows string `json:"binWindows,omitempty"`
		BinDarwin  string `json:"binDarwin,omitempty"`
		BinLinux   string `json:"binLinux,omitempty"`
		Bundle     string `json:"bundle,omitempty"`
		Config     string `json:"config,omitempty"`
		RepoUrl    string `json:"repoUrl"`
	}

	Credential struct {
		Service    string                 `json:"service"`
		Username   string                 `json:"username"`
		Credential map[string]interface{} `json:"credential"`
	}
	PermitMatcher struct {
		Pattern string   `yaml:"pattern"`
		Methods []string `yaml:"methods"`
	}

	DenyMatcher struct {
		Pattern      string              `yaml:"pattern"`
		RoleMappings map[string][]string `yaml:"roles"`
	}

	SecurityConstraints struct {
		Constraints       []DenyMatcher   `yaml:"constraints"`
		PublicConstraints []PermitMatcher `yaml:"publicConstraints"`
	}

	CredentialConfig struct {
		Field string `json:"field"`
		Type  string `json:"type"`
	}

	ConfigFile struct {
		CredentialConfig map[string][]CredentialConfig `json:"credentials"`
		CliVersionConfig CliVersionConfig              `json:"cliVersionPath"`
		RepositoryConfig []Repository                  `json:"repositories"`
		SPConfig         map[string]string             `json:"securityProviderConfig"`
	}

	CliVersionConfig struct {
		Url      string `json:"url"`
		Provider string `json:"provider"`
		Version  string `json:"cliversion"`
	}

	HealthEndpoints struct {
		VaultURL string
	}

	CreateUser struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	UserInfo struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	SecurityProviders struct {
		Providers map[string]SecurityManager
	}

	UserLogged struct {
		UserInfo UserInfo `json:"userInfo"`
		Roles    []string `json:"roles"`
		TTL      int64    `json:"ttl"`
		Org      string   `json:"org"`
	}
)

type Constraints interface {
	AuthorizationPath(bearerToken, path, method, org string) (bool, error)
	ValidatePublicConstraints(path, method string) bool
	ListRealmRoles(bearerToken, org string) ([]string, error)
}

type ConfigHealth interface {
	ReadHealthConfigs() map[string]HealthEndpoints
}

type ConfigCredential interface {
	ReadCredentialConfigs(org string) (map[string][]CredentialConfig, error)
}

type ConfigCliVersion interface {
	ReadCliVersionConfigs(org string) (CliVersionConfig, error)
}

type ConfigRepository interface {
	ReadRepositoryConfig(org string) ([]Repository, error)
}

type ConfigSecurityConstraints interface {
	ReadSecurityConstraints() SecurityConstraints
}

type Config interface {
	ConfigHealth
	ConfigCredential
	ConfigCliVersion
	ConfigRepository
	ConfigSecurityConstraints
}

type VaultManager interface {
	Write(key string, data map[string]interface{}) error
	Read(key string) (map[string]interface{}, error)
	List(key string) ([]interface{}, error)
	Delete(key string) error
	Start(*api.Client)
	Encrypt(data string) (string, error)
	Decrypt(data string) (string, error)
}

type VaultConfig interface {
	Start() (*api.Client, error)
}

type DefaultHandler interface {
	Handler() http.HandlerFunc
}

type CredentialHandler interface {
	HandleAdmin() http.HandlerFunc
	HandleMe() http.HandlerFunc
	HandleOrg() http.HandlerFunc
}

type MiddlewareHandler interface {
	Filter(next http.Handler) http.Handler
}

type ProviderHandler interface {
	TreeAllow(path, bToken, org string, repo Repository) (Tree, error)
	FilesFormulasAllow(path, bToken, org string, repo Repository) ([]byte, error)
	FindRepo(repos []Repository, repoName string) (Repository, error)
}

type SecurityManager interface {
	Login(username, password string) (User, LoginError)
	TTL() int64
}

type LoginError interface {
	Error() error
	Code() int
}

type User interface {
	Roles() []string
	UserInfo() UserInfo
}

type Configurator interface {
	LoadLoginHandler() DefaultHandler
	LoadCredentialConfigHandler() DefaultHandler
	LoadConfigHealth() DefaultHandler
	LoadUsageLoggerHandler() DefaultHandler
	LoadCliVersionHandler() DefaultHandler
	LoadRepositoryHandler() DefaultHandler
	LoadTreeHandler() DefaultHandler
	LoadFormulasHandler() DefaultHandler
	LoadMiddlewareHandler() MiddlewareHandler
	LoadCredentialHandler() CredentialHandler
	LoadHelloHandler() DefaultHandler
}

type WildcardPatternMatcher interface {
	Match() bool
}

type Slicer interface {
	Interface() ([]interface{}, error)
}
