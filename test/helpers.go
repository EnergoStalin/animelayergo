package test

import (
	"encoding/json"
	"os"
	"path"
)

type ConfigCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ResolveCredentialsFromDisk() (*ConfigCredentials, error) {
	bytes, err := os.ReadFile(path.Join(os.Getenv("HOME"), ".animelayer.json"))
	if err != nil {
		return nil, err
	}

	var value ConfigCredentials
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func getCreds() (string, string, error) {
	login, password := os.Getenv("ANIMELAYER_LOGIN"), os.Getenv("ANIMELAYER_PASSWORD")
	if login != "" && password != "" {
		return login, password, nil
	}

	creds, err := ResolveCredentialsFromDisk()
	if err != nil {
		return "", "", err
	}

	return creds.Login, creds.Password, nil
}
