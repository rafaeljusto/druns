package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/rafaeljusto/druns/core"
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

func (u URLs) GetHTTPS(name string) string {
	return u["baseHTTPS"] + u[name]
}

func (u URLs) GetHTTP(name string) string {
	return u["base"] + u[name]
}

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

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

////////////////////////////////////////////////
////////////////////////////////////////////////
////////////////////////////////////////////////

func EncryptPassword(password string) string {
	// TODO
	return password
}

func DecryptPassword(encryptedPassword string) string {
	// TODO
	return encryptedPassword
}
