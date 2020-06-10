package config

import (
	"fmt"
	"os"
	"reflect"
	"ritchie-server/server"
	"testing"
)

func TestConfiguration_ReadCliVersionConfigs(t *testing.T) {
	type fields struct {
		Configs             map[string]*server.ConfigFile
		SecurityConstraints server.SecurityConstraints
	}
	type args struct {
		org string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    server.CliVersionConfig
		outErr bool
	}{
		{
			name: "read cli version configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
						CliVersionConfig: server.CliVersionConfig{
							Url:      "http://localhost:8882/s3-version-mock",
							Provider: "s3",
						},
					},
				},
			},
			in: args{org: "zup"},
			out: server.CliVersionConfig{
				Url:      "http://localhost:8882/s3-version-mock",
				Provider: "s3",
			},
			outErr: false,
		},
		{
			name: "error read cli version configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
						CliVersionConfig: server.CliVersionConfig{
							Url:      "http://localhost:8882/s3-version-mock",
							Provider: "s3",
						},
					},
				},
			},
			in:     args{org: "error"},
			out:    server.CliVersionConfig{},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfiguration(tt.fields.Configs, tt.fields.SecurityConstraints)
			got, err := c.ReadCliVersionConfigs(tt.in.org)
			if (err != nil) != tt.outErr {
				t.Errorf("ReadCliVersionConfigs() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("ReadCliVersionConfigs() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestConfiguration_ReadCredentialConfigs(t *testing.T) {
	type fields struct {
		Configs             map[string]*server.ConfigFile
		SecurityConstraints server.SecurityConstraints
	}
	type args struct {
		org string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    map[string][]server.CredentialConfig
		outErr bool
	}{
		{
			name: "read credential configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
						CredentialConfig: map[string][]server.CredentialConfig{
							"credential1": {{Field: "Field", Type: "type"}},
							"credential2": {{Field: "field2", Type: "type"}},
						},
					},
				},
			},
			in: args{org: "zup"},
			out: map[string][]server.CredentialConfig{
				"credential1": {{Field: "Field", Type: "type"}},
				"credential2": {{Field: "field2", Type: "type"}},
			},
			outErr: false,
		},
		{
			name: "error read credential configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
						CredentialConfig: map[string][]server.CredentialConfig{
							"credential1": {{Field: "Field", Type: "type"}},
							"credential2": {{Field: "field2", Type: "type"}},
						},
					},
				},
			},
			in:     args{org: "error"},
			out:    nil,
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfiguration(tt.fields.Configs, tt.fields.SecurityConstraints)
			got, err := c.ReadCredentialConfigs(tt.in.org)
			if (err != nil) != tt.outErr {
				t.Errorf("ReadCredentialConfigs() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("ReadCredentialConfigs() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestConfiguration_ReadHealthConfigs(t *testing.T) {
	type fields struct {
		Configs             map[string]*server.ConfigFile
		SecurityConstraints server.SecurityConstraints
	}
	tests := []struct {
		name   string
		fields fields
		out    map[string]server.HealthEndpoints
	}{
		{
			name: "read health check configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
					},
				},
			},
			out: map[string]server.HealthEndpoints{
				"zup": {
					VaultURL: getVaultUrl(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfiguration(tt.fields.Configs, tt.fields.SecurityConstraints)
			if got := c.ReadHealthConfigs(); !reflect.DeepEqual(got, tt.out) {
				t.Errorf("ReadHealthConfigs() = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestConfiguration_ReadRepositoryConfig(t *testing.T) {
	type fields struct {
		Configs             map[string]*server.ConfigFile
		SecurityConstraints server.SecurityConstraints
	}
	type args struct {
		org string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    []server.Repository
		outErr bool
	}{
		{
			name: "read repository configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
						RepositoryConfig: []server.Repository{
							{
								Name:     "local",
								Priority: 0,
								TreePath: "path_whatever",
								Username: "",
								Password: "",
							},
							{
								Name:     "repository1",
								Priority: 1,
								TreePath: "path_whatever_repository1",
								Username: "optional",
								Password: "optional",
							},
						},
					},
				},
			},
			in: args{org: "zup"},
			out: []server.Repository{
				{
					Name:     "local",
					Priority: 0,
					TreePath: "path_whatever",
					Username: "",
					Password: "",
				},
				{
					Name:     "repository1",
					Priority: 1,
					TreePath: "path_whatever_repository1",
					Username: "optional",
					Password: "optional",
				},
			},
			outErr: false,
		},
		{
			name: "error read repository configuration",
			fields: fields{
				Configs: map[string]*server.ConfigFile{
					"zup": {
						RepositoryConfig: []server.Repository{
							{
								Name:     "local",
								Priority: 0,
								TreePath: "path_whatever",
								Username: "",
								Password: "",
							},
							{
								Name:     "repository1",
								Priority: 1,
								TreePath: "path_whatever_repository1",
								Username: "optional",
								Password: "optional",
							},
						},
					},
				},
			},
			in:     args{org: "error"},
			out:    nil,
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfiguration(tt.fields.Configs, tt.fields.SecurityConstraints)
			got, err := c.ReadRepositoryConfig(tt.in.org)
			if (err != nil) != tt.outErr {
				t.Errorf("ReadRepositoryConfig() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("ReadRepositoryConfig() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestConfiguration_ReadSecurityConstraints(t *testing.T) {
	type fields struct {
		Configs             map[string]*server.ConfigFile
		SecurityConstraints server.SecurityConstraints
	}
	tests := []struct {
		name   string
		fields fields
		out    server.SecurityConstraints
	}{
		{
			name: "read security constraints configuration",
			fields: fields{
				SecurityConstraints: server.SecurityConstraints{
					Constraints: []server.DenyMatcher{{
						Pattern:      "/test",
						RoleMappings: map[string][]string{"user": {"POST", "GET"}},
					}},
					PublicConstraints: []server.PermitMatcher{{
						Pattern: "/public",
						Methods: []string{"POST", "GET"},
					}},
				},
			},
			out: server.SecurityConstraints{
				Constraints: []server.DenyMatcher{{
					Pattern:      "/test",
					RoleMappings: map[string][]string{"user": {"POST", "GET"}},
				}},
				PublicConstraints: []server.PermitMatcher{{
					Pattern: "/public",
					Methods: []string{"POST", "GET"},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfiguration(tt.fields.Configs, tt.fields.SecurityConstraints)
			if got := c.ReadSecurityConstraints(); !reflect.DeepEqual(got, tt.out) {
				t.Errorf("ReadSecurityConstraints() = %v, out %v", got, tt.out)
			}
		})
	}
}

func getVaultUrl() string {
	p := "%s/sys/health"
	value := os.Getenv("VAULT_ADDR")
	if value == "" {
		value = "https://127.0.0.1:8200"
	}
	return fmt.Sprintf(p, value)
}
