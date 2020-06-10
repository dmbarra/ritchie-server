package security

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"ritchie-server/server"
	"ritchie-server/server/wpm"
)

type Authorization struct {
	config       server.Config
	vaultManager server.VaultManager
}

func NewAuthorization(c server.Config, v server.VaultManager) server.Constraints {
	return Authorization{
		config:       c,
		vaultManager: v,
	}
}

func (auth Authorization) AuthorizationPath(bearerToken, path, method, org string) (bool, error) {
	if org == "" {
		return false, fmt.Errorf("x-org not received. ")
	}
	roles, err := auth.ListRealmRoles(bearerToken, org)
	if err != nil {
		return false, err
	}
	return auth.validateConstraints(path, method, roles), nil
}

func (auth Authorization) ValidatePublicConstraints(path, method string) bool {

	sc := auth.config.ReadSecurityConstraints()

	for _, pc := range sc.PublicConstraints {
		if wpm.NewWildcardPattern(path, pc.Pattern).Match() {
			for _, m := range pc.Methods {
				if method == m {
					return true
				}
			}
		}
	}
	return false
}

func (auth Authorization) validateConstraints(path, method string, roles []string) bool {

	sc := auth.config.ReadSecurityConstraints()

	for _, pc := range sc.Constraints {
		if wpm.NewWildcardPattern(path, pc.Pattern).Match() {
			for _, role := range roles {
				rm := pc.RoleMappings[role]
				for _, m := range rm {
					if method == m {
						return true
					}
				}
			}
		}
	}
	return false
}

func (auth Authorization) ListRealmRoles(token, org string) ([]string, error) {
	if token != "" {
		t, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return nil, fmt.Errorf("failed decode token, error: %v", err)
		}
		tf, err := auth.vaultManager.Decrypt(string(t))
		if err != nil {
			return nil, errors.New("failed decrypt token")
		}
		var ul server.UserLogged
		err = json.Unmarshal([]byte(tf), &ul)
		if err != nil {
			return nil, errors.New("failed unmarshal token to user info")
		}
		if org != ul.Org {
			return nil, errors.New("receive org not equal token")
		}
		tokenTime := time.Unix(ul.TTL, 0)
		if time.Since(tokenTime).Seconds() > 0 {
			return nil, errors.New("token expired")
		}
		return ul.Roles, nil
	}
	return nil, errors.New("token is empty")
}
