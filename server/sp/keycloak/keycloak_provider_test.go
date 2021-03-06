package keycloak

import (
	"os"
	"reflect"
	"testing"
	"time"

	"ritchie-server/server"
)

func Test_keycloak_Login(t *testing.T) {
	type fields struct {
		config map[string]string
	}
	type args struct {
		username string
		password string
		totp     string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		outUser  server.User
		outError server.LoginError
	}{
		{
			name: "login success",
			fields: fields{
				config: dummyConfigKeycloak(""),
			},
			args: args{
				username: "user",
				password: "admin",
			},
			outUser: keycloakUser{
				roles: []string{"admin", "offline_access", "uma_authorization", "user"},
				userInfo: server.UserInfo{
					Name:     "user user",
					Username: "user",
					Email:    "user@user.com",
				},
			},
			outError: nil,
		},
		{
			name: "login failed",
			fields: fields{
				config: dummyConfigKeycloak(""),
			},
			args: args{
				username: "user",
				password: "failed",
			},
			outUser: nil,
			outError: keycloakError{
				code: 401,
				err:  nil,
			},
		},
		{
			name: "valid otp",
			fields: fields{
				config: map[string]string{
					"type":         "keycloak",
					"url":          "any url",
					"realm":        "ritchie",
					"clientId":     "user-login",
					"clientSecret": "user-login",
					"ttl":          "36000",
					"otp":          "true",
				},
			},
			args: args{
				username: "user",
				password: "failed",
			},
			outUser: nil,
			outError: keycloakError{
				code: 400,
				err:  ErrInvalidOpt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeycloakProvider(tt.fields.config)
			gotUser, gotError := k.Login(tt.args.username, tt.args.password, tt.args.totp)
			if gotUser != nil && tt.outUser != nil {
				if !reflect.DeepEqual(gotUser.UserInfo(), tt.outUser.UserInfo()) {
					t.Errorf("Login() gotUser.UserInfo() = %v, want %v", gotUser.UserInfo(), tt.outUser.UserInfo())
				}
				roles := make(map[string]string)
				for _, c := range gotUser.Roles() {
					roles[c] = c
				}
				for _, c := range tt.outUser.Roles() {
					if roles[c] == "" {
						t.Errorf("Error roles gotUser.Roles() = %v, want %v", gotUser.Roles(), tt.outUser.Roles())
					}
				}
			}
			if gotError != nil && tt.outError != nil {
				if gotError.Code() != tt.outError.Code() {
					t.Errorf("Login() gotError = %v, want %v", gotError, tt.outError)
				}
			}

		})
	}
}

func Test_keycloak_TTL(t *testing.T) {
	type fields struct {
		config map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "success",
			fields: fields{
				config: dummyConfigKeycloak(""),
			},
			want: 36000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeycloakProvider(tt.fields.config)
			ttl := k.TTL() - time.Now().Unix()
			if ttl != tt.want {
				t.Errorf("TTL() = %v, want %v", ttl, tt.want)
			}
		})
	}
}

func Test_keycloakConfig_Otp(t *testing.T) {
	type fields struct {
		config map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returned true",
			fields: fields{
				config: dummyConfigKeycloak("true"),
			},
			want: true,
		},
		{
			name: "returned false",
			fields: fields{
				config: dummyConfigKeycloak("false"),
			},
			want: false,
		},
		{
			name: "returned false (empty)",
			fields: fields{
				config: dummyConfigKeycloak(""),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeycloakProvider(tt.fields.config)
			if got := k.Otp(); got != tt.want {
				t.Errorf("Otp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func dummyConfigKeycloak(otp string) map[string]string {
	value := os.Getenv("KEYCLOAK_URL")
	if value == "" {
		value = "http://localhost:8080"
	}
	return map[string]string{
		"type":         "keycloak",
		"url":          value,
		"realm":        "ritchie",
		"clientId":     "user-login",
		"clientSecret": "user-login",
		"ttl":          "36000",
		"otp":          otp,
	}
}
