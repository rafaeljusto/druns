package protocol

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

var Translations map[string]Translation

type Translation map[msgCode]string

func init() {
	Translations = make(map[string]Translation)
}

func LoadTranslations(dirname string) error {
	Translations = map[string]Translation{}

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
		translation := make(map[msgCode]string)

		bytes, err := ioutil.ReadFile(path.Join(dirname, file.Name()))
		if err != nil {
			return errors.New(err)
		}

		if err := json.Unmarshal(bytes, &translation); err != nil {
			return errors.New(err)
		}

		Translations[language] = translation
	}

	return nil
}
