package credential

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"ritchie-server/server"
)

func (h Handler) HandleMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.processMeGet(w, r)
		case http.MethodPost:
			h.processMePost(w, r)
		default:
			http.NotFound(w, r)
		}
	}
}

func (h Handler) processMePost(w http.ResponseWriter, r *http.Request) {
	user, err := h.loadUser(*r)
	if err != nil {
		log.Error(fmt.Sprintf("Failed load user to request: %v", r))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	org := org(r)
	ctx := ctx(r)
	var c server.Credential
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Error("Failed to process request ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	c.Username = user.username

	if err := h.defaultValidate(c, org); len(err) > 0 {
		err := map[string]interface{}{"validationError": err}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if err := h.createCredential(org, ctx, c); err != nil {
		log.Error("Failed to create credential ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) processMeGet(w http.ResponseWriter, r *http.Request) {
	user, err := h.loadUser(*r)
	if err != nil {
		log.Error(fmt.Sprintf("Failed load user to request: %v", r))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	service := serviceFromPath(r.URL.Path)
	org := org(r)
	ctx := ctx(r)
	c := server.Credential{Service: service, Username: user.username}

	cre, err := h.findCredential(org, ctx, c)
	if err != nil {
		log.Error("Failed to retrieve credential ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if cre == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	c.Credential = cre
	data, _ := json.Marshal(c)
	w.Header().Add("Content-Type", "application/json")

	_, _ = w.Write(data)
}

func serviceFromPath(path string) string {
	return strings.Replace(path, "/credentials/me/", "", 1)
}
