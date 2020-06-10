package login

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"ritchie-server/server"
)

type Handler struct {
	securityProviders server.SecurityProviders
	vaultManager      server.VaultManager
}

func NewLoginHandler(s server.SecurityProviders, v server.VaultManager) server.DefaultHandler {
	return Handler{
		securityProviders: s,
		vaultManager:      v,
	}
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	Token string `json:"token"`
	TTL   int64  `json:"ttl"`
}

func (lh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var method = r.Method
		if http.MethodPost == method {
			lh.processPost(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (lh Handler) processPost(w http.ResponseWriter, r *http.Request) {
	var l login
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		log.Error("Failed to process Json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	organizationHeader := r.Header.Get(server.OrganizationHeader)
	if validErrs := l.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		er := json.NewEncoder(w).Encode(err)
		if er != nil {
			log.Printf("Error in Json Encode ")
			return
		}
		return
	}
	sp := lh.securityProviders.Providers[organizationHeader]
	if sp == nil {
		log.Printf("No provider security to org: %s", organizationHeader)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	lu, le := sp.Login(l.Username, l.Password)
	if le != nil {
		log.Printf("error login: %v", le.Error())
		w.WriteHeader(le.Code())
		return
	}
	resp := lh.createResponse(lu, organizationHeader, sp.TTL())
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Println("Error in Json Encode ")
		return
	}
}

func (lh Handler) createResponse(user server.User, org string, ttl int64) response {
	userLogged := server.UserLogged{
		UserInfo: user.UserInfo(),
		Roles:    user.Roles(),
		TTL:      ttl,
		Org:      org,
	}
	jb, _ := json.Marshal(userLogged)
	js := string(jb)
	cipher, err := lh.vaultManager.Encrypt(js)
	cipherB64 := base64.StdEncoding.EncodeToString([]byte(cipher))
	if err != nil {
		log.Fatal(err)
	}
	return response{
		Token: cipherB64,
		TTL:   ttl,
	}
}

func (l login) validate() url.Values {
	errs := url.Values{}
	if l.Username == "" {
		errs.Add("username", "The username field is required!")
	}
	if l.Password == "" {
		errs.Add("password", "The password field is required!")
	}
	return errs
}
