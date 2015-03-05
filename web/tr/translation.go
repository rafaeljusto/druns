package tr

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/rafaeljusto/druns/core/errors"
)

var (
	testLanguageFilename = regexp.MustCompile(`[a-z]{2,3}(-[A-Z]{2,3})?\.json`)
)

var translations map[string]translation

type translation map[errors.ValidationCode]string

func init() {
	translations = make(map[string]translation)
}

func LoadTranslations(dirname string) error {
	translations = map[string]translation{}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return errors.New(err)
	}

	for _, file := range files {
		if file.IsDir() || !testLanguageFilename.MatchString(file.Name()) {
			continue
		}

		language := file.Name()[:strings.Index(file.Name(), ".")]
		language = strings.Split(language, "-")[0]
		translation := make(map[errors.ValidationCode]string)

		bytes, err := ioutil.ReadFile(path.Join(dirname, file.Name()))
		if err != nil {
			return errors.New(err)
		}

		if err := json.Unmarshal(bytes, &translation); err != nil {
			return errors.New(err)
		}

		translations[language] = translation
	}

	return nil
}
