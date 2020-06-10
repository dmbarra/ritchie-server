package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestConfigHealth_Handler(t *testing.T) {
	type fields struct {
		Config server.Config
		Path   string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name:   "success",
			fields: fields{Config: mock.DummyConfig(), Path: "/health"},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {

					orgs := []org{{
						NameOrg: "zup",
						Healths: healthStruct{
							Services: []service{
								{
									ServiceType: "VAULT",
									Health:      "UP",
								},
							},
							Status: "UP",
						},
					},
					}

					js, _ := json.Marshal(orgs)
					_, err := w.Write(js)
					if err != nil {
						fmt.Sprintln("Error in Write ")
						return
					}
				}
			}(),
		},
		{
			name:   "not found",
			fields: fields{Config: mock.DummyConfig("not"), Path: "/nothealth"},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {

					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := NewConfigHealth(tt.fields.Config)

			r, _ := http.NewRequest(http.MethodGet, tt.fields.Path, bytes.NewReader([]byte{}))

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}

			if g.Code == 200 {
				if !reflect.DeepEqual(g.Body, w.Body) {
					t.Errorf("Handler returned wrong body: got %v \n want %v", g.Body, w.Body)
				}
			}
		})
	}
}