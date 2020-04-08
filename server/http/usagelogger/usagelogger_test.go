package usagelogger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		Config  server.Config
		method  string
		org     string
		command interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		{
			name: "success",
			fields: fields{
				Config: mock.DummyConfig(),
				method: http.MethodPost,
				org:    "zup",
				command: cmdUser{
					Username: "user",
					Cmd:      "rit sample",
				},
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusOK)
				}
			}(),
		},
		{
			name: "not found",
			fields: fields{
				Config: mock.DummyConfig(),
				method: http.MethodGet,
				org:    "zup",
				command: cmdUser{
					Username: "user",
					Cmd:      "rit sample",
				},
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
		{
			name: "bad request",
			fields: fields{
				Config:  mock.DummyConfig(),
				method:  http.MethodPost,
				org:     "zup",
				command: cmdUser{},
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusBadRequest)
				}
			}(),
		},
		{
			name: "error",
			fields: fields{
				Config:  mock.DummyConfig(),
				method:  http.MethodPost,
				org:     "zup",
				command: "test",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "org not found",
			fields: fields{
				Config: mock.DummyConfig(),
				method: http.MethodPost,
				org:    "notfound",
				command: cmdUser{
					Username: "user",
					Cmd:      "rit sample",
				},
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mu := NewUsageLoggerHandler(tt.fields.Config)

			payloadJson, _ := json.Marshal(tt.fields.command)

			r, _ := http.NewRequest(tt.fields.method, "/metrics/use", bytes.NewReader(payloadJson))

			r.Header.Add(server.OrganizationHeader, tt.fields.org)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			tt.want.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			mu.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}
		})
	}
}
