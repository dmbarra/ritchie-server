package ul

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"ritchie-server/server"
)

type Handler struct {
}

type cmdUser struct {
	Username string `json:"username"`
	Cmd      string `json:"command"`
}

func NewUsageLoggerHandler() server.DefaultHandler {
	return Handler{}
}

func (mu Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodPost:
			mu.processPost(w, r)
		default:
			http.NotFound(w, r)
		}
	}
}

func (mu Handler) processPost(w http.ResponseWriter, r *http.Request) {
	var mUser cmdUser
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&mUser)

	if err != nil {
		log.Error("Failed to process Json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	if validErrs := mUser.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	mUserJSON, _ := json.Marshal(mUser)
	log.Info(string(mUserJSON))
	w.WriteHeader(http.StatusOK)

}

func (cmd cmdUser) validate() url.Values {
	errs := url.Values{}
	if cmd.Username == "" {
		errs.Add("username", "The username field is required!")
	}
	if cmd.Cmd == "" {
		errs.Add("command", "The command field is required!")
	}
	return errs
}
