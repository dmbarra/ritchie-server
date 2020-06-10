package login

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		s       server.SecurityProviders
		v       server.VaultManager
		method  string
		org     string
		payload string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name: "success",
			fields: fields{
				s: server.SecurityProviders{
					Providers: map[string]server.SecurityManager{
						"zup": mock.SecurityManagerMock{
							U: mock.UserMock{
								R: []string{"rit_admin"},
								U: server.UserInfo{
									Name:     "User User",
									Username: "user",
									Email:    "user@user.com",
								},
							},
							L: nil,
							T: 34000,
						},
					},
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodPost,
				org:     "zup",
				payload: `{"username": "user", "password":"admin"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
				}
			}(),
		},
		{
			name: "sp not found",
			fields: fields{
				s: server.SecurityProviders{
					Providers: map[string]server.SecurityManager{
						"zup": mock.SecurityManagerMock{
							U: mock.UserMock{
								R: []string{"rit_admin"},
								U: server.UserInfo{
									Name:     "User User",
									Username: "user",
									Email:    "user@user.com",
								},
							},
							L: nil,
							T: 34000,
						},
					},
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodPost,
				org:     "not",
				payload: `{"username": "user", "password":"admin"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
		{
			name: "failed login",
			fields: fields{
				s: server.SecurityProviders{
					Providers: map[string]server.SecurityManager{
						"zup": mock.SecurityManagerMock{
							U: nil,
							L: mock.LoginErrorMock{
								E: errors.New("error"),
								C: 401,
							},
							T: 34000,
						},
					},
				},
				method:  http.MethodPost,
				org:     "zup",
				payload: `{"username": "user", "password":"failed"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}(),
		},
		{
			name: "method not found",
			fields: fields{
				method: http.MethodGet,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Header().Set("Content-type", "text/plain; charset=utf-8")
				}
			}(),
		},
		{
			name: "failed input",
			fields: fields{
				method:  http.MethodPost,
				org:     "zup",
				payload: `1`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "empty fields",
			fields: fields{
				method:  http.MethodPost,
				org:     "zup",
				payload: `{}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-type", "application/json")
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewLoginHandler(tt.fields.s, tt.fields.v)
			var b []byte
			if len(tt.fields.payload) > 0 {
				b = append(b, []byte(tt.fields.payload)...)
			}
			r, _ := http.NewRequest(tt.fields.method, "/test", bytes.NewReader(b))
			r.Header.Add(server.OrganizationHeader, tt.fields.org)

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}

			if g.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
				t.Errorf("Wrong content type. Got %v want %v", g.Header().Get("Content-Type"), w.Header().Get("Content-Type"))
			}
		})
	}
}
