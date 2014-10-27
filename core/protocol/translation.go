package protocol

import (
	"br/core"
	"io/ioutil"
	"jsonconf"
	"path"
	"regexp"
	"strings"
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
		return core.NewError(err)
	}

	for _, file := range files {
		if file.IsDir() || !testLanguageFilename.MatchString(file.Name()) {
			continue
		}

		language := file.Name()[:strings.Index(file.Name(), ".")]
		language = strings.Split(language, "-")[0]
		translation := make(map[msgCode]string)

		err := jsonconf.LoadConfigFile(path.Join(dirname, file.Name()), &translation)
		if err != nil {
			return core.NewError(err)
		}

		Translations[language] = translation
	}

	return nil
}
