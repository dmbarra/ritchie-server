package fph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"ritchie-server/server"
)

const (
	providerHttp = "HTTP"
	providerS3   = "S3"
)

type Handler struct {
	authorization server.Constraints
}

func NewProviderHandler(a server.Constraints) server.ProviderHandler {
	return Handler{
		authorization: a,
	}
}

func (hp Handler) TreeAllow(path, bToken, org string, repo server.Repository) (server.Tree, error) {
	rTree, err := treeRemote(path, repo)
	if err != nil {
		return rTree, err
	}
	roles, err := hp.authorization.ListRealmRoles(bToken, org)
	if err != nil {
		return rTree, err
	}
	rfind := make(map[string]interface{})
	for _, r := range roles {
		rfind[strings.ToUpper(r)] = r
	}
	ft := server.Tree{}
	ft.Version = rTree.Version
	for _, c := range rTree.Commands {
		if len(c.Roles) > 0 {
			for _, r := range c.Roles {
				if rfind[strings.ToUpper(r)] != nil {
					ft.Commands = append(ft.Commands, c)
				}
			}
		} else {
			ft.Commands = append(ft.Commands, c)
		}
	}
	if repo.ReplaceRepoUrl != "" {
		for _, c := range ft.Commands {
			if c.Formula != nil {
				if c.Formula.RepoUrl != "" {
					c.Formula.RepoUrl = repo.ReplaceRepoUrl
				}
			}
		}
	}
	return ft, nil
}

func (hp Handler) FilesFormulasAllow(path, bToken, org string, repo server.Repository) ([]byte, error) {
	tr, err := hp.TreeAllow(repo.TreePath, bToken, org, repo)
	if err != nil {
		return nil, err
	}
	roles, err := hp.authorization.ListRealmRoles(bToken, org)
	if err != nil {
		return nil, err
	}

	rfind := make(map[string]interface{})
	for _, r := range roles {
		rfind[strings.ToUpper(r)] = r
	}
	p := strings.Replace(path, "/formulas/", "", 1)
	s := strings.Split(p, "/")
	key := strings.ReplaceAll(p, "/"+s[len(s)-1], "")
	for _, c := range tr.Commands {
		if c.Formula != nil {
			if c.Formula.Path == key {
				switch repo.Provider.Type {
				case providerHttp:
					return processHttp(path, repo)
				case providerS3:
					return processS3(path, repo)
				}
			}
		}
	}
	return nil, nil
}

func (hp Handler) FindRepo(repos []server.Repository, repoName string) (server.Repository, error) {
	var repository server.Repository
	for _, r := range repos {
		if r.Name == repoName {
			repository = r
			break
		}
	}
	if repository.Name == "" {
		return repository, fmt.Errorf("No repo with name %s\n", repoName)
	}
	return repository, nil
}

func processS3(path string, repo server.Repository) ([]byte, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(repo.Provider.Region)},
	)
	if err != nil {
		return nil, err
	}
	buf := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	s3obj := s3.GetObjectInput{
		Bucket: aws.String(repo.Provider.Bucket),
		Key:    aws.String(path),
	}
	_, err = downloader.Download(buf,
		&s3obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func processHttp(path string, repo server.Repository) ([]byte, error) {
	url := fmt.Sprintf("%s%s", repo.Provider.Remote, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func treeRemote(tPath string, repo server.Repository) (server.Tree, error) {
	switch repo.Provider.Type {
	case providerHttp:
		return loadTreeFileHttp(tPath, repo)
	case providerS3:
		return loadTreeFileS3(tPath, repo)
	default:
		return server.Tree{}, fmt.Errorf("provider %q, not valid. Verify our repo config. Repo name: %q", repo.Provider.Type, repo.Name)
	}
}

func loadTreeFileHttp(path string, repo server.Repository) (server.Tree, error) {
	url := fmt.Sprintf("%s%s", repo.Provider.Remote, path)
	var response server.Tree
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		return response, fmt.Errorf("%d - failed to get index for %s\n", resp.StatusCode, url)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return response, err
	}
	return response, nil
}

func loadTreeFileS3(path string, repo server.Repository) (server.Tree, error) {
	var response server.Tree
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(repo.Provider.Region)},
	)
	if err != nil {
		return response, err
	}
	buf := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(repo.Provider.Bucket),
			Key:    aws.String(path),
		})
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		return response, err
	}
	return response, nil
}
