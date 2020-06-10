package vault

import (
	"encoding/base64"
	"fmt"
	"os"

	vlogin "github.com/ZupIT/go-vault-session/pkg/login"
	"github.com/ZupIT/go-vault-session/pkg/token"
	"github.com/hashicorp/vault/api"

	"ritchie-server/server"
	"ritchie-server/server/slicer"

	log "github.com/sirupsen/logrus"
)

const vaultPath = "ritchie/credential/%s"
const vaultEncrypt = "ritchie/transit/encrypt/%s"
const vaultDecrypt = "ritchie/transit/decrypt/%s"
const ritKey = "ritchie_key"

type Manager struct {
	client *api.Client
}

func NewVaultManager(c *api.Client) server.VaultManager {
	return Manager{client: c}
}

func (vm Manager) Write(key string, data map[string]interface{}) error {
	vm.setToken()
	if _, err := vm.client.Logical().Write(pathResolver(key), data); err != nil {
		log.Error("Vault write error ", err)
		return err
	}
	return nil
}

func (vm Manager) Read(key string) (map[string]interface{}, error) {
	vm.setToken()
	res, err := vm.client.Logical().Read(pathResolver(key))

	if err != nil {
		log.Error("Vault read error ", err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return res.Data, nil
}

func (vm Manager) Delete(key string) error {
	vm.setToken()
	_, err := vm.client.Logical().Delete(pathResolver(key))
	if err != nil {
		log.Error("Vault delete error ", err)
		return err
	}

	return nil
}

func (vm Manager) List(key string) ([]interface{}, error) {
	vm.setToken()
	res, err := vm.client.Logical().List(pathResolver(key))
	if err != nil {
		log.Error("Vault list error ", err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	a := res.Data["keys"]

	keys, err := slicer.NewSlicer(a).Interface()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (vm Manager) Start(c *api.Client) {
	l := vlogin.NewHandler(c)
	s := l.Handle()

	ch := make(chan string)
	r := token.NewHandler(c, s, ch)

	r.Handle()
	go func() {
		for {
			msg := <-ch
			log.Info(msg)
		}
	}()

}

func (vm Manager) Encrypt(data string) (string, error) {
	vm.setToken()
	path := fmt.Sprintf(vaultEncrypt, ritKey)
	body := make(map[string]interface{})
	body["plaintext"] = base64.StdEncoding.EncodeToString([]byte(data))
	res, err := vm.client.Logical().Write(path, body)
	if err != nil {
		log.Println("Vault encrypt error", err)
		return "", err
	}
	return res.Data["ciphertext"].(string), nil
}

func (vm Manager) Decrypt(data string) (string, error) {
	vm.setToken()
	path := fmt.Sprintf(vaultDecrypt, ritKey)
	body := make(map[string]interface{})
	body["ciphertext"] = data
	res, err := vm.client.Logical().Write(path, body)
	if err != nil {
		log.Println("Vault decrypt error", err)
		return "", err
	}
	plainText, err := base64.StdEncoding.DecodeString(res.Data["plaintext"].(string))
	if err != nil {
		log.Println("Vault decode error", err)
		return "", err
	}
	return string(plainText), nil
}

func (vm Manager) setToken() {
	vm.client.SetToken(os.Getenv(api.EnvVaultToken))
}

func pathResolver(key string) string {
	return fmt.Sprintf(vaultPath, key)
}
