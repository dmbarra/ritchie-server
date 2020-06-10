package starter

import (
	"fmt"
	"reflect"
	"ritchie-server/server"
	"testing"
)

func TestConfigurator_LoadLoginHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "login.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}

			if got := reflect.TypeOf(c.LoadLoginHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadLoginHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadCliVersionHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "cliversion.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadCliVersionHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadCliVersionHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadConfigHealth(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "health.ConfigHealth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadConfigHealth()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadConfigHealth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadCredentialConfigHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "credential.ConfigHandler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadCredentialConfigHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadCredentialConfigHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadCredentialHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "credential.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadCredentialHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadCredentialHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadUsageLoggerHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "ul.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadUsageLoggerHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadUsageLoggerHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadMiddlewareHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "middleware.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadMiddlewareHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadMiddlewareHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadRepositoryHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "repository.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadRepositoryHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadRepositoryHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadTreeHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "tree.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadTreeHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadTreeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadFormulasHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "formulas.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadFormulasHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadFormulasHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigurator_LoadHelloHandler(t *testing.T) {
	type fields struct {
		conf         server.Config
		vaultManager server.VaultManager
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "correct type",
			fields: fields{},
			want:   "hello.Handler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configurator{
				conf:         tt.fields.conf,
				vaultManager: tt.fields.vaultManager,
			}
			if got := reflect.TypeOf(c.LoadHelloHandler()); fmt.Sprint(got) != tt.want {
				t.Errorf("LoadHelloHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "run as intended",
			want:    "starter.Configurator",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfiguration()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := reflect.TypeOf(got); fmt.Sprint(got) != tt.want {
				t.Errorf("NewConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
