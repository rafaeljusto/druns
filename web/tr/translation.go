package tr

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/rafaeljusto/druns/core"
)

var (
	testLanguageFilename = regexp.MustCompile(`[a-z]{2,3}(-[A-Z]{2,3})?\.json`)
)

var translations map[string]translation

type translation map[core.ValidationErrorCode]string

func init() {
	translations = make(map[string]translation)
}

func LoadTranslations(dirname string) error {
	translations = map[string]translation{}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return core.NewError(err)
	}

	for _, file := range files {
		if file.IsDir() || !testLanguageFilename.MatchString(file.Name()) {
			continue
		}

		language := file.Name()[:strings.Index(file.Name(), ".")]
		language = strings.Split(language, "-")[0]
		translation := make(map[core.ValidationErrorCode]string)

		bytes, err := ioutil.ReadFile(path.Join(dirname, file.Name()))
		if err != nil {
			return core.NewError(err)
		}

		if err := json.Unmarshal(bytes, &translation); err != nil {
			return core.NewError(err)
		}

		translations[language] = translation
	}

	return nil
}
