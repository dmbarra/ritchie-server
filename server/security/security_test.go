package security

import (
	"encoding/json"
	"errors"
	"log"
	"testing"
	"time"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestAuthorization_AuthorizationPath(t *testing.T) {
	type fields struct {
		c server.Config
		v server.VaultManager
	}
	type args struct {
		token  string
		path   string
		method string
		org    string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    bool
		outErr bool
	}{
		{
			name:   "empty org",
			fields: fields{},
			in:     args{},
			out:    false,
			outErr: true,
		},
		{
			name: "empty token",
			in: args{
				token: "",
				org:   "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "invalid token",
			in: args{
				token: "invalid",
				org:   "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "decrypt error",
			fields: fields{
				c: mock.DummyConfig(),
				v: mock.VaultMock{
					Err:     errors.New("error"),
					ErrList: nil,
					Keys:    nil,
					Data:    "",
				},
			},
			in: args{
				token: "dG9rZW4=",
				org:   "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "failed unmarshal token",
			fields: fields{
				c: mock.DummyConfig(),
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
					Data:    "failed unmarshal",
				},
			},
			in: args{
				token: "dG9rZW4=",
				org:   "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "org diff token",
			fields: fields{
				c: mock.DummyConfig(),
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
					Data:    jsonUserLogged(3600),
				},
			},
			in: args{
				token: "dG9rZW4=",
				org:   "fail",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "token expired",
			fields: fields{
				c: mock.DummyConfig(),
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
					Data:    jsonUserLogged(-1),
				},
			},
			in: args{
				token: "dG9rZW4=",
				org:   "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "Validate constrains",
			fields: fields{
				c: mock.DummyConfig(),
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
					Data:    jsonUserLogged(36000),
				},
			},
			in: args{
				token:  "dG9rZW4=",
				org:    "zup",
				method: "GET",
				path:   "/validate",
			},
			out:    true,
			outErr: false,
		},
		{
			name: "Invalid constrains",
			fields: fields{
				c: mock.DummyConfig(),
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
					Data:    jsonUserLogged(36000),
				},
			},
			in: args{
				token:  "dG9rZW4=",
				org:    "zup",
				method: "GET",
				path:   "/test",
			},
			out:    false,
			outErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthorization(tt.fields.c, tt.fields.v)
			got, err := auth.AuthorizationPath(tt.in.token, tt.in.path, tt.in.method, tt.in.org)
			if (err != nil) != tt.outErr {
				t.Errorf("AuthorizationPath() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if got != tt.out {
				t.Errorf("AuthorizationPath() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func jsonUserLogged(ttl int64) string {
	t := time.Now().Unix() + ttl
	u := server.UserLogged{
		UserInfo: server.UserInfo{
			Name:     "test",
			Username: "test",
			Email:    "test@test.com",
		},
		Roles: []string{"rit_user", "admin"},
		TTL:   t,
		Org:   "zup",
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Fatal("error json.Marshal(u)")
	}
	return string(b)
}

func TestAuthorization_ValidatePublicConstraints(t *testing.T) {
	type fields struct {
		c server.Config
		v server.VaultManager
	}
	type args struct {
		path   string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		out    bool
	}{
		{
			name: "Validate constrains",
			fields: fields{
				c: mock.DummyConfig(),
			},
			args: args{
				method: "GET",
				path:   "/public",
			},
			out: true,
		},
		{
			name: "Invalid constrains",
			fields: fields{
				c: mock.DummyConfig(),
			},
			args: args{
				method: "GET",
				path:   "/invalid",
			},
			out: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthorization(tt.fields.c, tt.fields.v)
			if got := auth.ValidatePublicConstraints(tt.args.path, tt.args.method); got != tt.out {
				t.Errorf("ValidatePublicConstraints() = %v, out %v", got, tt.out)
			}
		})
	}
}
