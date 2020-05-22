package ph

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestHandler_FindRepo(t *testing.T) {
	type fields struct {
		authorization server.Constraints
	}
	type args struct {
		repos         []server.Repository
		repoName      string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    server.Repository
		outErr bool
	}{
		{
			name: "find commons",
			in: args{
				repos:    mock.DummyRepoList(),
				repoName: "commons",
			},
			out:    mock.DummyRepoList()[0],
			outErr: false,
		},
		{
			name: "repo not found",
			in: args{
				repos:    mock.DummyRepoList(),
				repoName: "notfound",
			},
			out:    server.Repository{},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := Handler{
				authorization: tt.fields.authorization,
			}
			got, err := hp.FindRepo(tt.in.repos, tt.in.repoName)
			if (err != nil) != tt.outErr {
				t.Errorf("FindRepo() error = %v, wantErr %v", err, tt.outErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("FindRepo() got = %v, want %v", got, tt.out)
			}
		})
	}
}

func TestHandler_TreeAllow(t *testing.T) {
	type fields struct {
		authorization server.Constraints
	}
	type args struct {
		path   string
		bToken string
		org    string
		repo   server.Repository
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    server.Tree
		outErr bool
	}{
		{
			name: "response tree with allow commands",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
			},
			in: args{
				path:   "/tree/tree.json",
				bToken: "",
				org:    "",
				repo:   mock.DummyRepo(),
			},
			out:    treeRoleUser(),
			outErr: false,
		},
		{
			name: "provider error",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
			},
			in: args{
				path:   "/tree/tree.json",
				bToken: "",
				org:    "",
				repo:   mock.DummyRepo("ERROR"),
			},
			out:    server.Tree{},
			outErr: true,
		},
		{
			name: "list roles error",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: false,
					E: errors.New("error"),
					R: nil,
				},
			},
			in: args{
				bToken: "",
				org:    "",
				path:   "/tree/tree.json",
				repo:   mock.DummyRepo(),
			},
			out:    server.Tree{},
			outErr: true,
		},
		{
			name: "tree not found",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
			},
			in: args{
				bToken: "",
				org:    "",
				path:   "/tree/tree-notfound.json",
				repo:   mock.DummyRepo(),
			},
			out:    server.Tree{},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := Handler{
				authorization: tt.fields.authorization,
			}
			got, err := hp.TreeAllow(tt.in.path, tt.in.bToken, tt.in.org, tt.in.repo)
			if (err != nil) != tt.outErr {
				t.Errorf("TreeAllow() error = %v, wantErr %v", err, tt.outErr)
				return
			}
			commands := make(map[string]*server.Command)
			for _, c := range got.Commands {
				commands[c.Parent+c.Usage] = &c
			}
			for _, c := range tt.out.Commands {
				if commands[c.Parent+c.Usage] == nil {
					t.Errorf("Commands receive in tree error gotT = %v, outT %v", got, tt.out)
				}
			}
		})
	}
}

func TestHandler_FilesFormulasAllow(t *testing.T) {
	type fields struct {
		authorization server.Constraints
	}
	type args struct {
		path   string
		bToken string
		org    string
		repo   server.Repository
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    bool
		outErr bool
	}{
		{
			name: "allow",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
			},
			in: args{
				path:   "/formulas/aws/terraform/config.json",
				bToken: "",
				org:    "",
				repo:   mock.DummyRepo(),
			},
			out:    true,
			outErr: false,
		},
		{
			name: "not allow",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"NO_RULE"},
				},
			},
			in: args{
				path:   "/formulas/aws/terraform/config.json",
				bToken: "",
				org:    "",
				repo:   mock.DummyRepo(),
			},
			out:    false,
			outErr: false,
		},
		{
			name: "load roles error",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: errors.New("error"),
					R: []string{},
				},
			},
			in: args{
				path:   "/formulas/aws/terraform/config.json",
				bToken: "",
				org:    "",
				repo:   mock.DummyRepo(),
			},
			out:    false,
			outErr: true,
		},
		{
			name: "allow formula without role",
			fields: fields{
				authorization: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
			},
			in: args{
				path:   "/formulas/scaffold/coffee-go/config.json",
				bToken: "",
				org:    "",
				repo:   mock.DummyRepo(),
			},
			out:    true,
			outErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := Handler{
				authorization: tt.fields.authorization,
			}
			got, err := hp.FilesFormulasAllow(tt.in.path, tt.in.bToken, tt.in.org, tt.in.repo)
			if (err != nil) != tt.outErr {
				t.Errorf("FilesFormulasAllow() error = %v, wantErr %v", err, tt.outErr)
				return
			}
			if len(got) > 0 != tt.out {
				t.Errorf("FilesFormulasAllow() got = %v, want %v", got, tt.out)
			}
		})
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
		log.Fatal("error Unmarshal tree")
	}
	return s
}

