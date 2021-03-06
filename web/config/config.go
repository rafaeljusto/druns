package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core/errors"
)

var (
	DrunsConfig Config
)

func init() {
	DrunsConfig.URLs = make(map[string]string)
	DrunsConfig.Files = make(map[string][]string)
}

type Config struct {
	Home string

	Server struct {
		IP   string
		Port int
		TLS  struct {
			PrivKey string
			Cert    string
		}
	}

	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}

	Log struct {
		Host string
		Port int
	}

	Paths struct {
		ProtocolTranslations string
		WebTranslations      string
		HTMLTemplates        string
		WebAssets            string
	}

	URLs      URLs
	Files     map[string][]string
	Languages []string

	Session struct {
		Secret  string
		Timeout Duration
	}

	ClassValue float64
}

func (c Config) TLS() (string, string) {
	return path.Join(c.Home, c.Server.TLS.Cert),
		path.Join(c.Home, c.Server.TLS.PrivKey)
}

func (c Config) HTMLTemplates(language, handlerName string) []string {
	p := fmt.Sprintf(c.Paths.HTMLTemplates, language)

	templates := make([]string, len(c.Files[handlerName]))
	copy(templates, c.Files[handlerName])

	for i, template := range templates {
		templates[i] = path.Join(c.Home, p, template)
	}

	return templates
}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

type URLs map[string]string

func (u URLs) GetHTTPS(name string, params ...string) string {
	path := u["baseHTTPS"] + u[name]

	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	return path
}

func (u URLs) GetHTTP(name string, params ...string) string {
	path := u["base"] + u[name]

	if len(params) > 0 {
		path += "?" + strings.Join(params, "&")
	}

	return path
}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

func LoadConfig(path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(err)
	}

	if err := json.Unmarshal(bytes, &DrunsConfig); err != nil {
		return errors.New(err)
	}

	return nil
}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	duration, err := time.ParseDuration(strings.Trim(string(data), "\""))
	if err != nil {
		return err
	}

	d.Duration = duration
	return nil
}
