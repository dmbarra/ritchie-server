package ldap

import (
	"reflect"
	"testing"
	"time"

	"ritchie-server/server"
)

func Test_ldapConfig_Login(t *testing.T) {
	type fields struct {
		config map[string]string
	}
	type args struct {
		username string
		password string
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
				config: map[string]string{
					"type":               "ldap",
					"base":               "dc=example,dc=org",
					"host":               "localhost",
					"serverName":         "ldap.example.org",
					"port":               "389",
					"useSSL":             "false",
					"skipTLS":            "true",
					"insecureSkipVerify": "false",
					"bindDN":             "cn=admin,dc=example,dc=org",
					"bindPassword":       "admin",
					"userFilter":         "(uid=%s)",
					"groupFilter":        "(memberUid=%s)",
					"attributeUsername":  "uid",
					"attributeName":      "givenName",
					"attributeEmail":     "mail",
					"ttl":                "36000",
				},
			},
			args: args{
				username: "user",
				password: "user",
			},
			outUser: ldapUser{
				roles: []string{"rit_user"},
				userInfo: server.UserInfo{
					Name:     "Test",
					Username: "user",
					Email:    "test@test.com.br",
				},
			},
			outError: nil,
		},
		{
			name: "login failed",
			fields: fields{
				config: map[string]string{
					"type":               "ldap",
					"base":               "dc=example,dc=org",
					"host":               "localhost",
					"serverName":         "ldap.example.org",
					"port":               "389",
					"useSSL":             "false",
					"skipTLS":            "true",
					"insecureSkipVerify": "false",
					"bindDN":             "cn=admin,dc=example,dc=org",
					"bindPassword":       "admin",
					"userFilter":         "(uid=%s)",
					"groupFilter":        "(memberUid=%s)",
					"attributeUsername":  "uid",
					"attributeName":      "givenName",
					"attributeEmail":     "mail",
					"ttl":                "36000",
				},
			},
			args: args{
				username: "user",
				password: "failed",
			},
			outUser: nil,
			outError: ldapError{
				code: 401,
				err:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewLdapProvider(tt.fields.config)
			gotUser, gotError := k.Login(tt.args.username, tt.args.password)
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

func Test_ldapConfig_TTL(t *testing.T) {
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
				config: map[string]string{
					"type":               "ldap",
					"base":               "dc=example,dc=org",
					"host":               "localhost",
					"serverName":         "ldap.example.org",
					"port":               "389",
					"useSSL":             "false",
					"skipTLS":            "true",
					"insecureSkipVerify": "false",
					"bindDN":             "cn=admin,dc=example,dc=org",
					"bindPassword":       "admin",
					"userFilter":         "(uid=%s)",
					"groupFilter":        "(memberUid=%s)",
					"attributeUsername":  "uid",
					"attributeName":      "givenName",
					"attributeEmail":     "mail",
					"ttl":                "36000",
				},
			},
			want: 36000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewLdapProvider(tt.fields.config)
			ttl := k.TTL() - time.Now().Unix()
			if ttl != tt.want {
				t.Errorf("TTL() = %v, want %v", ttl, tt.want)
			}
		})
	}
}
