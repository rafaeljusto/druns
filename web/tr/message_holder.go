package tr

import (
	"fmt"

	"github.com/rafaeljusto/druns/core"
	"github.com/rafaeljusto/druns/core/log"
	"github.com/rafaeljusto/druns/web/config"
)

type MessageHolder struct {
	language string
}

func NewMessageHolder(language string) MessageHolder {
	return MessageHolder{language}
}

func (m MessageHolder) Get(c core.ValidationErrorCode, args ...interface{}) string {
	choosenLanguage := m.language

	translation, exists := translations[choosenLanguage]
	if !exists {
		// Choose another language
		for _, language := range config.DrunsConfig.Languages {
			if translation, exists = translations[language]; exists {
				choosenLanguage = language
				break
			}
		}

		log.Warningf("Language “%s” not found, using language “%s” instead",
			m.language, choosenLanguage)
	}

	msg, exists := translation[c]
	if !exists {
		// Choose another translation
		for _, language := range config.DrunsConfig.Languages {
			if translation, exists = translations[language]; exists {
				if msg, exists = translation[c]; exists {
					choosenLanguage = language
					break
				}
			}
		}

		// Configuration error, code missing from JSON file
		log.Warningf("Language “%s” doesn't has a translation for code “%s”, "+
			"using language “%s” instead", m.language, c, choosenLanguage)
	}

	return fmt.Sprintf(msg, args...)
}
