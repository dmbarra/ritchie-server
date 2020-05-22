package tree

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/config"
	"ritchie-server/server/mock"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		config    server.Config
		auth      server.Constraints
		providerH server.ProviderHandler
		path      string
		method    string
		org       string
		repoName  string
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		{
			name: "tree allow user",
			fields: fields{
				config: mock.DummyConfig(),
				auth: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				providerH: mock.ProviderHandlerMock{
					T: treeRoleUser(),
					R: server.Repository{
						Name:           "commons",
						Priority:       0,
						TreePath:       "/tree/tree.json",
						ServerUrl:      "http://localhost:3000",
						ReplaceRepoUrl: "http://localhost:3000/formulas",
						Provider:       server.Provider{
							Type:   "HTTP",
							Remote: "http://localhost:8882",
						},
					},
				},
				path:     "/tree/tree.json",
				method:   http.MethodGet,
				org:      "zup",
				repoName: "commons",
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
					err := json.NewEncoder(w).Encode(treeRoleUser())
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
				path:     "/tree/tree.json",
				method:   http.MethodPost,
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
				path:     "/tree/tree.json",
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
				path:     "/tree/tree.json",
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
				path:     "/tree/tree.json",
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
			name: "internal error",
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
				path:     "/tree/tree.json",
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

			if g.Code == http.StatusOK {
				var got server.Tree
				if err := json.Unmarshal(g.Body.Bytes(), &got); err != nil {
					log.Fatal("Error Unmarshal server.Tree")
				}
				var out server.Tree
				if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
					log.Fatal("Error Unmarshal server.Tree")
				}
				commands := make(map[string]*server.Command)
				for _, c := range got.Commands {
					commands[c.Parent+c.Usage] = &c
				}
				for _, c := range out.Commands {
					if commands[c.Parent+c.Usage] == nil {
						t.Errorf("Commands receive in tree error gotT = %v, outT %v", got, out)
					}
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

func treeRoleUser() server.Tree {
	js := `{
        "commands": [
          {
            "usage": "aws",
            "help": "Apply Aws objects",
            "parent": "root",
            "roles" : ["USER"]
          },
          {
            "usage": "apply",
            "help": "Apply Aws objects",
            "parent": "root_aws",
            "roles" : ["USER"]
          },
          {
            "usage": "terraform",
            "help": "Apply Aws terraform objects",
            "formula": {
              "path": "aws/terraform",
              "bin": "terraform-cli-${so}",
              "bundle": "${so}.zip",
              "repoUrl": "https://commons-repo.ritchiecli.io/formulas"
            },
            "parent": "root_aws_apply",
            "roles" : ["USER"]
          },
          {
            "usage": "scaffold",
            "help": "Manipulate scaffold objects",
            "parent": "root"
          },
          {
            "usage": "generate",
            "help": "Generates a scaffold by some template",
            "parent": "root_scaffold"
          },
          {
            "usage": "coffee-go",
            "help": "Generates a project by coffee template in Go",
            "formula": {
              "path": "scaffold/coffee-go",
              "bin": "coffee-go-${so}",
              "bundle": "${so}.zip",
              "repoUrl": "https://commons-repo.ritchiecli.io/formulas"
            },
            "parent": "root_scaffold_generate"
          }
        ],
        "version": "1.0.0"
      }`
	var s server.Tree
	if err := json.Unmarshal([]byte(js), &s); err != nil {
		log.Fatal(err)
	}
	return s
}
