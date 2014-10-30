package tr

import (
	"io/ioutil"
	"jsonconf"
	"path"
	"regexp"
	"strings"
)

var (
	testLanguageFilename = regexp.MustCompile(`[a-z]{2,3}(-[A-Z]{2,3})?\.json`)
)

var translations map[string]translation

type translation map[Code]string

func init() {
	translations = make(map[string]translation)
}

func LoadTranslations(dirname string) error {
	translations = map[string]translation{}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || !testLanguageFilename.MatchString(file.Name()) {
			continue
		}

		language := file.Name()[:strings.Index(file.Name(), ".")]
		language = strings.Split(language, "-")[0]
		translation := make(map[Code]string)

		err := jsonconf.LoadConfigFile(path.Join(dirname, file.Name()), &translation)
		if err != nil {
			return err
		}

		translations[language] = translation
	}

	return nil
}
