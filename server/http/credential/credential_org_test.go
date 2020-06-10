package credential

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestHandler_HandleOrg(t *testing.T) {
	type fields struct {
		v       server.VaultManager
		c       server.Config
		method  string
		org     string
		ctx     string
		payload string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name:   "method not allowed",
			fields: fields{method: http.MethodGet},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Allow", http.MethodPost)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}(),
		},
		{
			name: "post credential success",
			fields: fields{method: http.MethodPost, v: mock.VaultMock{
				ReturnMap:  map[string]interface{}{"a": "b"},
			},
				payload: mock.DummyCredentialAdmin(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
				}
			}(),
		},
		{
			name: "post credential invalid json",
			fields: fields{method: http.MethodPost, v: mock.VaultMock{
				ReturnMap:  map[string]interface{}{"a": "b"},
			},
				payload: "failed",
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "post credential bad request",
			fields: fields{method: http.MethodPost, v: mock.VaultMock{
				ReturnMap:  map[string]interface{}{"a": "b"},
			},
				payload: mock.DummyCredentialBadRequest(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
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
			fields: fields{method: http.MethodPost, v: mock.VaultMock{
				Err:      errors.New("error"),
				ReturnMap:  map[string]interface{}{"a": "b"},
			},
				payload: mock.DummyCredentialAdmin(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "post credential empty service",
			fields: fields{method: http.MethodPost, v: mock.VaultMock{
				Err:      errors.New("error"),
				ReturnMap:  map[string]interface{}{"a": "b"},
			},
				payload: mock.DummyCredentialEmpty(),
				ctx:     "default",
				org:     "zup",
				c:       mock.DummyConfig(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Add("Content-Type", "application/json")
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
			r.Header.Add(server.AuthorizationHeader, "dGVzdA==")
			r.Header.Add(server.ContextHeader, tt.fields.ctx)
			r.Header.Add(server.OrganizationHeader, tt.fields.org)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.HandleOrg().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v out %v", g.Code, w.Code)
			}

			if g.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
				t.Errorf("Wrong content type. Got %v out %v", g.Header().Get("Content-Type"), w.Header().Get("Content-Type"))
			}
		})
	}
}
