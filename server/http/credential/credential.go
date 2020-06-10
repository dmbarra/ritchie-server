package credential

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"ritchie-server/server"
)

const (
	credentialVaultPath    = "/%s/%s/%s"
	orgCredentialVaultPath = "/%s/%s"
)

type user struct {
	username   string
	email      string
	name       string
}

type Handler struct {
	v server.VaultManager
	c server.Config
}

func NewCredentialHandler(v server.VaultManager, c server.Config) server.CredentialHandler {
	return Handler{v: v, c: c}
}

func (h Handler) defaultValidate(c server.Credential, org server.Org) map[string]string {
	errs := make(map[string]string)

	credentials, _ := h.c.ReadCredentialConfigs(string(org))

	found := false

	for k := range credentials {
		if k == c.Service {
			found = true
			break
		}
	}

	if !found {
		errs["credential"] = "The credential is not valid."
	}
	if c.Service == "" {
		errs["service"] = "The service field is required."
	}
	if len(c.Credential) == 0 {
		errs["credential"] = "The credential field must not be empty."
	}
	return errs
}

func (h Handler) findCredential(org server.Org, ctx server.Ctx, c server.Credential) (map[string]interface{}, error) {
	path := fmt.Sprintf(credentialVaultPath, ctxResolver(org, ctx), c.Username, c.Service)
	credential, err := h.v.Read(path)
	if err != nil {
		return nil, err
	}

	// By default, we search credentials by the user,
	// but if the user doesn't have a credential,
	// we will search by the organization.
	if credential == nil {
		credential, err = h.findOrgCredential(org, ctx, c)
		if err != nil {
			return nil, err
		}
	}

	return credential, nil
}

func (h Handler) findOrgCredential(org server.Org, ctx server.Ctx, c server.Credential) (map[string]interface{}, error) {
	path := fmt.Sprintf(orgCredentialVaultPath, ctxResolver(org, ctx), c.Service)
	credential, err := h.v.Read(path)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (h Handler) createCredential(org server.Org, ctx server.Ctx, c server.Credential) error {
	path := fmt.Sprintf(credentialVaultPath, ctxResolver(org, ctx), c.Username, c.Service)
	if err := h.v.Write(path, c.Credential); err != nil {
		return err
	}
	return nil
}

func ctxResolver(org server.Org, ctx server.Ctx) string {
	if ctx != "" {
		return fmt.Sprintf("%s_%s", org, ctx)
	}

	return string(org)
}

func org(r *http.Request) server.Org {
	return server.Org(r.Header.Get(server.OrganizationHeader))
}

func ctx(r *http.Request) server.Ctx {
	return server.Ctx(r.Header.Get(server.ContextHeader))
}

func (h Handler) loadUser(r http.Request) (user, error) {
	authorizationToken := r.Header.Get(server.AuthorizationHeader)
	t, err := base64.StdEncoding.DecodeString(authorizationToken)
	if err != nil {
		return user{}, fmt.Errorf("failed decode token, error: %v", err)
	}
	tf, err := h.v.Decrypt(string(t))
	if err != nil {
		return user{}, errors.New("failed decrypt token")
	}
	var ul server.UserLogged
	err = json.Unmarshal([]byte(tf), &ul)
	if err != nil {
		return user{}, errors.New("failed unmarshal token to user info")
	}
	return user{
		username: ul.UserInfo.Username,
		email:    ul.UserInfo.Email,
		name:     ul.UserInfo.Name,
	}, nil
}
