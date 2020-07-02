package httpmock

import (
	"net/http"
	"net/http/httptest"
	"strings"

	log "github.com/sirupsen/logrus"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func LoadServerHttp() *httptest.Server {
	return httptest.NewUnstartedServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if strings.Contains(req.URL.String(), "/tree/tree.json") {
				loadResponse(res, "../mock/json/tree.json")
			}
			if strings.Contains(req.URL.String(), "/formulas/aws/terraform/config.json") {
				loadResponse(res, "../mock/json/awsterraformconfig.json")
			}
			if strings.Contains(req.URL.String(), "/formulas/scaffold/coffee-go/config.json") {
				loadResponse(res, "../mock/json/scaffoldcoffee_goconfig.json")
			}
		}),
	)
}

func loadResponse(res http.ResponseWriter, file string) {
	res.WriteHeader(200)
	res.Header().Set("Content-Type", "application/json")
	_, err := res.Write(mock.LoadJson(file))
	if err != nil {
		log.Error(err)
	}

}

func GenerateRepoWithMock(ts *httptest.Server) server.Repository {
	return server.Repository{
		Name:           "commons",
		Priority:       0,
		TreePath:       "/tree/tree.json",
		ServerUrl:      ts.URL,
		ReplaceRepoUrl: ts.URL + "/formulas",
		Provider: server.Provider{
			Type:   "HTTP",
			Remote: ts.URL,
		},
	}
}