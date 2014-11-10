package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/rafaeljusto/druns/core"
)

var (
	DrunsConfig Config
)

func init() {
	DrunsConfig.Files = make(map[string][]string)
}

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

	Paths struct {
		Home                 string
		ProtocolTranslations string
		WebTranslations      string
		HTMLTemplates        string
		WebAssets            string
	}

	Files     map[string][]string
	Languages []string
}

func (c Config) HTMLTemplates(language, handlerName string) []string {
	p := fmt.Sprintf(c.Paths.HTMLTemplates, language)

	templates := make([]string, len(c.Files[handlerName]))
	copy(templates, c.Files[handlerName])

	for i, template := range templates {
		templates[i] = path.Join(c.Paths.Home, p, template)
	}

	return templates
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
