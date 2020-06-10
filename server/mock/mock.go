package mock

import (
	"os"

	"github.com/hashicorp/vault/api"

	"ritchie-server/server"
	"ritchie-server/server/config"
)

const (
	keycloakUrl   = "KEYCLOAK_URL"
	cliVersionUrl = "CLI_VERSION_URL"
	remoteUrl     = "REMOTE_URL"
)

func DummyConfig(args ...string) server.Config {
	return config.Configuration{
		Configs: DummyConfigMap(args...),
		SecurityConstraints: server.SecurityConstraints{
			Constraints: []server.DenyMatcher{{
				Pattern:      "/validate",
				RoleMappings: map[string][]string{"admin": {"POST", "GET"}},
			}},
			PublicConstraints: []server.PermitMatcher{{
				Pattern: "/public",
				Methods: []string{"POST", "GET"},
			}},
		},
	}
}

func DummyConfigMap(args ...string) map[string]*server.ConfigFile {
	keycloakUrl := getEnv(keycloakUrl, "http://localhost:8080")
	remoteUrl := getEnv(remoteUrl, "http://localhost:8882")
	realm := "ritchie"
	clientId := "user-login"
	clientSecret := "user-login"
	if len(args) > 0 && args[0] != "" {
		keycloakUrl = args[0]
	}
	if len(args) > 1 && args[1] != "" {
		realm = args[1]
	}
	if len(args) > 2 && args[2] != "" {
		clientId = args[2]
	}
	if len(args) > 3 && args[3] != "" {
		clientSecret = args[3]
	}
	return map[string]*server.ConfigFile{
		"zup": {
			SPConfig: map[string]string{
				"type":         "keycloak",
				"url":          keycloakUrl,
				"realm":        realm,
				"clientId":     clientId,
				"clientSecret": clientSecret,
			},
			CredentialConfig: map[string][]server.CredentialConfig{
				"credential1": {{Field: "Field", Type: "type"}},
				"credential2": {{Field: "field2", Type: "type"}},
			},
			CliVersionConfig: server.CliVersionConfig{
				Url:      getEnv(cliVersionUrl, "http://localhost:8882/s3-version-mock"),
				Provider: "s3",
			},
			RepositoryConfig: []server.Repository{
				{
					Name:           "commons",
					Priority:       0,
					TreePath:       "/tree/tree.json",
					ServerUrl:      "http://localhost:3000",
					ReplaceRepoUrl: "http://localhost:3000/formulas",
					Provider: server.Provider{
						Type:   "HTTP",
						Remote: remoteUrl,
					},
				},
				{
					Name:           "test1",
					Priority:       1,
					TreePath:       "/tree/tree-test1.json",
					ServerUrl:      "http://localhost:3000",
					ReplaceRepoUrl: "http://localhost:3000/formulas",
					Provider: server.Provider{
						Type:   "HTTP",
						Remote: remoteUrl,
					},
				},
				{
					Name:           "test-repo",
					Priority:       2,
					TreePath:       "/tree/tree-test2.json",
					ServerUrl:      "http://localhost:3000",
					ReplaceRepoUrl: "http://localhost:3000/formulas",
					Provider: server.Provider{
						Type:   "S3",
						Bucket: "local",
						Region: "sa-east-1",
					},
				},
			},
		}}
}

// Cli Version
func DummyConfigCliVersionUrlNotFound() server.Config {
	return config.Configuration{
		Configs: map[string]*server.ConfigFile{
			"zup": {
				CliVersionConfig: server.CliVersionConfig{
					Provider: "s3",
				},
			}},
	}
}
func DummyConfigCliVersionUrlWrong() server.Config {
	return config.Configuration{
		Configs: map[string]*server.ConfigFile{
			"zup": {
				CliVersionConfig: server.CliVersionConfig{
					Url:      "wrong",
					Provider: "s3",
				},
			}},
	}
}
func DummySecurityConstraints() server.SecurityConstraints {
	return server.SecurityConstraints{
		Constraints: []server.DenyMatcher{{
			Pattern:      "/test",
			RoleMappings: map[string][]string{"user": {"POST", "GET"}},
		}},
		PublicConstraints: []server.PermitMatcher{{
			Pattern: "/public",
			Methods: []string{"POST", "GET"},
		}},
	}
}

// Credential
func DummyCredential() string {
	return `{
	"service": "credential1",
		"credential": {
			"username": "test",
			"token": "token"
		}
	}`
}
func DummyCredentialEmpty() string {
	return `{
	"username": "Ubijara",
	"service": "",
		"credential": {
		}
	}`
}
func DummyCredentialAdmin() string {
	return `{
	"username": "Ubijara",
	"service": "credential1",
		"credential": {
			"username": "test",
			"token": "token"
		}
	}`
}
func DummyCredentialBadRequest() string {
	return `{
	"service": "invalid",
		"credential": {
			"username": "test",
			"token": "token"
		}
	}`
}

func DummyRepo(args ...string) server.Repository {
	remote := getEnv(remoteUrl, "http://localhost:8882")
	tp := "HTTP"
	if len(args) > 0 {
		tp = args[0]
	}
	return server.Repository{
		Name:           "commons",
		Priority:       0,
		TreePath:       "/tree/tree.json",
		ServerUrl:      "http://localhost:3000",
		ReplaceRepoUrl: "http://localhost:3000/formulas",
		Provider: server.Provider{
			Type:   tp,
			Remote: remote,
		},
	}
}

func DummyRepoList() []server.Repository {
	remote := getEnv(remoteUrl, "http://localhost:8882")
	return []server.Repository{
		{
			Name:           "commons",
			Priority:       0,
			TreePath:       "/tree/tree.json",
			ServerUrl:      "http://localhost:3000",
			ReplaceRepoUrl: "http://localhost:3000/formulas",
			Provider: server.Provider{
				Type:   "HTTP",
				Remote: remote,
			},
		},
		{
			Name:           "test1",
			Priority:       1,
			TreePath:       "/tree/tree-test1.json",
			ServerUrl:      "http://localhost:3000",
			ReplaceRepoUrl: "http://localhost:3000/formulas",
			Provider: server.Provider{
				Type:   "HTTP",
				Remote: remote,
			},
		},
		{
			Name:           "test2",
			Priority:       2,
			TreePath:       "/tree/tree-test2.json",
			ServerUrl:      "http://localhost:3000",
			ReplaceRepoUrl: "http://localhost:3000/formulas",
			Provider: server.Provider{
				Type:   "HTTP",
				Remote: remote,
			},
		},
	}
}

// server.SecurityManager mock
type SecurityManagerMock struct {
	U server.User
	L server.LoginError
	T int64
}

func (s SecurityManagerMock) Login(username, password string) (server.User, server.LoginError) {
	return s.U, s.L
}
func (s SecurityManagerMock) TTL() int64 {
	return s.T
}

type LoginErrorMock struct {
	E error
	C int
}

func (le LoginErrorMock) Error() error {
	return le.E
}

func (le LoginErrorMock) Code() int {
	return le.C
}

type UserMock struct {
	R []string
	U server.UserInfo
}

func (u UserMock) Roles() []string {
	return u.R
}

func (u UserMock) UserInfo() server.UserInfo {
	return u.U
}

// server.ValtManager mock
type VaultMock struct {
	Err        error
	ErrList    error
	Keys       []interface{}
	Data       string
	ReturnMap  map[string]interface{}
	ErrDecrypt error
}

func (v VaultMock) Write(string, map[string]interface{}) error {
	return v.Err
}
func (v VaultMock) Read(string) (map[string]interface{}, error) {
	return v.ReturnMap, v.Err
}
func (v VaultMock) List(string) ([]interface{}, error) {
	return v.Keys, v.ErrList
}
func (v VaultMock) Delete(string) error {
	return v.Err
}
func (v VaultMock) Start(*api.Client) {
}

func (v VaultMock) Encrypt(data string) (string, error) {
	return v.Data, nil
}
func (v VaultMock) Decrypt(data string) (string, error) {
	return v.Data, v.ErrDecrypt
}

type AuthorizationMock struct {
	B bool
	E error
	R []string
}

func (d AuthorizationMock) AuthorizationPath(bearerToken, path, method, org string) (bool, error) {
	return d.B, d.E
}
func (d AuthorizationMock) ValidatePublicConstraints(path, method string) bool {
	return d.B
}
func (d AuthorizationMock) ListRealmRoles(bearerToken, org string) ([]string, error) {
	if d.E != nil {
		return nil, d.E
	}
	var new []string
	new = append(new, d.R...)
	return new, d.E
}

type ProviderHandlerMock struct {
	T  server.Tree
	B  []byte
	R  server.Repository
	ER error
	ET error
}

func (ph ProviderHandlerMock) TreeAllow(path, bToken, org string, repo server.Repository) (server.Tree, error) {
	return ph.T, ph.ET
}
func (ph ProviderHandlerMock) FilesFormulasAllow(path, bToken, org string, repo server.Repository) ([]byte, error) {
	return ph.B, ph.ET
}
func (ph ProviderHandlerMock) FindRepo(repos []server.Repository, repoName string) (server.Repository, error) {
	return ph.R, ph.ER
}

func getEnv(key, def string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return def
}
