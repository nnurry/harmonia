package config

import "golang.org/x/crypto/ssh"

type SSH struct {
	User string
	Host string
	Port int

	HostKeyCallback ssh.HostKeyCallback

	PasswordAuth   PasswordAuthSSH
	PrivateKeyAuth PrivateKeyAuthSSH
}

type PasswordAuthSSH struct {
	Password string
}

type PrivateKeyAuthSSH struct {
	PrivateKeyPath string
	Passphrase     string
}
