package formulas

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/config"
	"ritchie-server/server/mock"
)

// Config type that represents formula configDummy
type configDummy struct {
	Name        string  `json:"name"`
	Command     string  `json:"command"`
	Description string  `json:"description"`
	Language    string  `json:"language"`
	Inputs      []input `json:"inputs"`
}

// Input type that represents input configDummy
type input struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Default string   `json:"default"`
	Label   string   `json:"label"`
	Items   []string `json:"items"`
	Cache   cache    `json:"cache"`
}

type cache struct {
	Active   bool   `json:"active"`
	Qtd      int    `json:"qtd"`
	NewLabel string `json:"newLabel"`
}

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		config    server.Config
		auth      server.Constraints
		providerH server.ProviderHandler
		method    string
		path      string
		org       string
		repoName  string
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		{
			name: "allow config",
			fields: fields{
				config: mock.DummyConfig(),
				auth: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				providerH: mock.ProviderHandlerMock{
					R: server.Repository{
						Name:           "commons",
						Priority:       0,
						TreePath:       "/tree/tree.json",
						ServerUrl:      "http://localhost:3000",
						ReplaceRepoUrl: "http://localhost:3000/formulas",
						Provider: server.Provider{
							Type:   "HTTP",
							Remote: "http://localhost:8882",
						},
					},
					B: configJsonWantByte(),
				},
				method:   http.MethodGet,
				path:     "/formulas/aws/terraform/config.json",
				org:      "zup",
				repoName: "commons",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
					err := json.NewEncoder(w).Encode(configJsonWant())
					if err != nil {
						fmt.Sprintln("Error in Encode Json ")
						return
					}
				}
			}(),
		},
		{
			name: "not found post",
			fields: fields{
				config: mock.DummyConfig(),
				auth: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				method:   http.MethodPost,
				path:     "/formulas/aws/terraform/config.json",
				org:      "zup",
				repoName: "commons",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
		{
			name: "empty repo",
			fields: fields{
				config: configEmptyRepo(),
				auth: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				path:     "/formulas/aws/terraform/config.json",
				method:   http.MethodGet,
				org:      "zup",
				repoName: "commons",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
		{
			name: "no org",
			fields: fields{
				config: configEmptyRepo(),
				auth: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				path:     "/formulas/aws/terraform/config.json",
				method:   http.MethodGet,
				org:      "no",
				repoName: "commons",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
		{
			name: "repo name not found",
			fields: fields{
				config: mock.DummyConfig(),
				auth: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				providerH: mock.ProviderHandlerMock{
					ER: errors.New("error"),
				},
				path:     "/formulas/aws/terraform/config.json",
				method:   http.MethodGet,
				org:      "zup",
				repoName: "not found",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
		{
			name: "internal error load files",
			fields: fields{
				config: mock.DummyConfig(),
				auth: mock.AuthorizationMock{
					B: true,
					E: errors.New("error"),
					R: []string{"USER"},
				},
				providerH: mock.ProviderHandlerMock{
					ET: errors.New("error"),
				},
				path:     "/formulas/aws/terraform/config.json",
				method:   http.MethodGet,
				org:      "zup",
				repoName: "commons",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := NewConfigHandler(tt.fields.config, tt.fields.auth, tt.fields.providerH)

			r, _ := http.NewRequest(tt.fields.method, tt.fields.path, bytes.NewReader([]byte{}))

			r.Header.Add(server.OrganizationHeader, tt.fields.org)
			r.Header.Add(repoNameHeader, tt.fields.repoName)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			tt.want.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			mu.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}

			if g.Code == 200 {
				var got configDummy
				if err := json.Unmarshal(g.Body.Bytes(), &got); err != nil {
					log.Fatal("Error unmarshal configDummy")
				}
				var out configDummy
				if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
					log.Fatal("Error unmarshal configDummy")
				}
				if !reflect.DeepEqual(got, out) {
					t.Errorf("Handler returned wrong body: got %v \n want %v", g.Body, w.Body)
				}
			}
		})
	}
}

func configEmptyRepo() server.Config {
	return config.Configuration{
		Configs: map[string]*server.ConfigFile{
			"zup": {},
		},
	}
}

func configJsonWantByte() []byte {
	return []byte( `{
        "description": "Apply terraform on AWS",
        "inputs": [
          {
            "name": "repository",
            "type": "text",
            "default": "https://github.com/zup/terraform-mock",
            "label": "Select your repository URL: ",
            "items": [
              "https://github.com/zup/mock-1",
              "https://github.com/zup/mock-2"
            ]
          },
          {
            "name": "terraform_path",
            "type": "text",
            "default": "src",
            "label": "Type your terraform files path [src]: "
          },
          {
            "name": "environment",
            "type": "text",
            "label": "Type your environment name: [ex.: qa, prod]"
          },
          {
            "name": "git_user",
            "type": "CREDENTIAL_GITHUB_USERNAME"
          },
          {
            "name": "git_token",
            "type": "CREDENTIAL_GITHUB_TOKEN"
          },
          {
            "name": "aws_access_key_id",
            "type": "CREDENTIAL_AWS_ACCESSKEYID"
          },
          {
            "name": "aws_secret_access_key",
            "type": "CREDENTIAL_AWS_SECRETACCESSKEY"
          }
        ]
      }`)
}

func configJsonWant() configDummy {
	var c configDummy
	if err := json.Unmarshal(configJsonWantByte(), &c); err != nil {
		log.Fatal("Error Unmarshal")
	}
	return c
}
