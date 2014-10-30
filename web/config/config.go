package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rafaeljusto/druns/core"
)

var (
	DrunsConfig Config
)

type Config struct {
	Server struct {
		IP   string
		Port int
	}

	Database struct {
		Name string
		URI  string
	}

	Log struct {
		Host string
		Port int
	}

	Languages []string
}

func LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return core.NewError(err)
	}

	if err := json.Unmarshal(bytes, &DrunsConfig); err != nil {
		return core.NewError(err)
	}

	return nil
}
