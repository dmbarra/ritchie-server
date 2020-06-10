package credential

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestHandler_HandlerMe(t *testing.T) {
	type fields struct {
		v       server.VaultManager
		c       server.Config
		method  string
		org     string
		ctx     string
		payload string
		auth    string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name: "method not found",
			fields: fields{
				method: http.MethodPatch,
				auth:   "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.NotFound(w, r)
				}
			}(),
		},
		{
			name: "get credential success",
			fields: fields{
				method: http.MethodGet,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				auth: "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
				}
			}(),
		},
		{
			name: "get error decode",
			fields: fields{
				method: http.MethodGet,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				auth: "error8",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}(),
		},
		{
			name: "get error decrypt",
			fields: fields{
				method: http.MethodGet,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
					ErrDecrypt: errors.New("error"),
				},
				auth: "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}(),
		},
		{
			name: "get error unmarshal",
			fields: fields{
				method: http.MethodGet,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      "err",
				},
				auth: "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}(),
		},
		{
			name: "post credential success",
			fields: fields{
				method: http.MethodPost,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				payload: mock.DummyCredential(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
				auth:    "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
				}
			}(),
		},
		{
			name: "post error decrypt",
			fields: fields{
				method: http.MethodPost,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				payload: mock.DummyCredential(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
				auth:    "test8",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
				}
			}(),
		},
		{
			name: "post credential invalid json",
			fields: fields{
				method: http.MethodPost,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				payload: "failed",
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
				auth:    "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "post credential bad request",
			fields: fields{
				method: http.MethodPost,
				v: mock.VaultMock{
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				payload: mock.DummyCredentialBadRequest(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
				auth:    "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Add("Content-Type", "application/json")
				}
			}(),
		},
		{
			name: "post credential error write",
			fields: fields{
				method: http.MethodPost,
				v: mock.VaultMock{
					Err:       errors.New("error"),
					ReturnMap: map[string]interface{}{"a": "b"},
					Data:      userLoggedJson(),
				},
				payload: mock.DummyCredential(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
				auth:    "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "get credential error",
			fields: fields{
				method: http.MethodGet,
				v: mock.VaultMock{
					Err:  errors.New("error"),
					Data: userLoggedJson(),
				},
				auth: "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "get credential not found",
			fields: fields{
				method: http.MethodGet,
				v:      mock.VaultMock{Data: userLoggedJson()},
				auth:   "dGVzdA==",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewCredentialHandler(tt.fields.v, tt.fields.c)
			var b []byte
			if len(tt.fields.payload) > 0 {
				b = append(b, []byte(tt.fields.payload)...)
			}
			r, _ := http.NewRequest(tt.fields.method, "/test", bytes.NewReader(b))
			r.Header.Add(server.AuthorizationHeader, tt.fields.auth)
			r.Header.Add(server.ContextHeader, tt.fields.ctx)
			r.Header.Add(server.OrganizationHeader, tt.fields.org)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.HandleMe().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v out %v", g.Code, w.Code)
			}

			if g.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
				t.Errorf("Wrong content type. Got %v out %v", g.Header().Get("Content-Type"), w.Header().Get("Content-Type"))
			}
		})
	}
}

func userLoggedJson() string {
	u := server.UserLogged{
		UserInfo: server.UserInfo{
			Name:     "test",
			Username: "test",
			Email:    "test@test.com",
		},
		Roles: []string{"rit_user", "rit_admin"},
		TTL:   0,
		Org:   "zup",
	}
	b, err := json.Marshal(u)
	if err != nil {
		log.Fatal("error json.Marshal(u)")
	}
	return string(b)
}
