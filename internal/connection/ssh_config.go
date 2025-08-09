package connection

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	User string `json:"user"`
	Host string `json:"host"`
	Port int    `json:"port"`

	HostKeyCallbackName string `json:"hostkey_callback_name"`

	PasswordAuth   passwordAuthSSHConfig   `json:"password_auth_config"`
	PrivateKeyAuth privateKeyAuthSSHConfig `json:"privkey_auth_config"`
}

type passwordAuthSSHConfig struct {
	Password string `json:"password"`
}

type privateKeyAuthSSHConfig struct {
	PrivateKeyPath string `json:"path"`
	Passphrase     string `json:"passphrase"`
}

func (cfg SSHConfig) HostKeyCallback(callbackName string) (ssh.HostKeyCallback, error) {
	switch callbackName {
	case "InsecureIgnoreHostKey":
		return ssh.InsecureIgnoreHostKey(), nil
	}
	return nil, errors.New("unsupported host key callback")
}

func (cfg SSHConfig) ParsePasswordAuth() (ssh.AuthMethod, error) {
	if cfg.PasswordAuth.Password == "" {
		return nil, fmt.Errorf("empty password")
	}
	return ssh.Password(cfg.PasswordAuth.Password), nil
}

func (cfg SSHConfig) ParsePrivateKeyAuth() (ssh.AuthMethod, error) {
	if cfg.PrivateKeyAuth.PrivateKeyPath == "" {
		return nil, fmt.Errorf("empty private key path")
	}

	keyContent, err := os.ReadFile(cfg.PrivateKeyAuth.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("can't read private key: %v", err)
	}

	var signer ssh.Signer

	if cfg.PrivateKeyAuth.Passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyContent, []byte(cfg.PrivateKeyAuth.Passphrase))
		if err != nil {
			return nil, fmt.Errorf("can't parse private key with passphrase: %v", err)
		}
	} else {
		signer, err = ssh.ParsePrivateKey(keyContent)
		if err != nil {
			return nil, fmt.Errorf("can't parse private key: %v", err)
		}
	}

	return ssh.PublicKeys(signer), nil
}
