package health

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ritchie-server/server"
)

type org struct {
	NameOrg string       `json:"nameOrg"`
	Healths healthStruct `json:"healths"`
}
type healthStruct struct {
	Services []service `json:"services"`
	Status   string    `json:"status"`
}
type service struct {
	ServiceType string `json:"type"`
	Health      string `json:"health"`
}

type ConfigHealth struct {
	Config server.Config
}

func NewConfigHealth(configuration server.Config) server.DefaultHandler {
	return ConfigHealth{Config: configuration}
}

func (ch ConfigHealth) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/health" {
			http.NotFound(w, r)
		} else {
			orgs := []org{}
			var vaultCheck bool
			var err error

			for key, h := range ch.Config.ReadHealthConfigs() {
				if err != nil {
					log.Error("Keycloak health checking error ", err)
				}
				vaultCheck, err = healthCheckUrl(h.VaultURL)
				if err != nil {
					log.Error("Vault health checking error ", err)
				}
				orgX := org{key, healthStruct{Services: []service{
					{ServiceType: "VAULT", Health: healthCheckStatus(vaultCheck)},
				}, Status: healthCheckStatus(vaultCheck)}}
				orgs = append(orgs, orgX)
			}

			resp, _ := json.Marshal(orgs)
			if !vaultCheck {
				w.WriteHeader(http.StatusInternalServerError)
			}
			fmt.Fprint(w, string(resp))
		}
	})
}

func healthCheckStatus(value bool) string {
	if value {
		return "UP"
	}
	return "DOWN"
}

func healthCheckUrl(url string) (bool, error) {

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}
