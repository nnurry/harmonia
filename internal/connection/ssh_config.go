package connection

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	User string
	Host string
	Port int

	HostKeyCallback ssh.HostKeyCallback

	PasswordAuth   passwordAuthSSHConfig
	PrivateKeyAuth privateKeyAuthSSHConfig
}

type passwordAuthSSHConfig struct {
	Password string
}

type privateKeyAuthSSHConfig struct {
	PrivateKeyPath string
	Passphrase     string
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
