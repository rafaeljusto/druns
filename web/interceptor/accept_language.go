package interceptor

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/web/config"
	"github.com/rafaeljusto/druns/web/tr"
)

type languageMessager interface {
	SetLanguage(language string)
	SetMessages(messageHolder tr.MessageHolder)
}

type AcceptLanguage struct {
	trama.NopAJAXInterceptor
	handler languageMessager
}

func NewAcceptLanguage(h languageMessager) *AcceptLanguage {
	return &AcceptLanguage{handler: h}
}

func (i AcceptLanguage) Before(w http.ResponseWriter, r *http.Request) {
	acceptLanguage := r.Header.Get("Accept-Language")
	acceptLanguageParts := strings.Split(acceptLanguage, ",")

	var selectedLanguage string
	var selectedQuality float64

	for _, part := range acceptLanguageParts {
		languageAndOptions := strings.Split(part, ";")

		language := languageAndOptions[0]
		var quality float64 = 1 // By default is quatility 100%

		for i := 1; i < len(languageAndOptions); i++ {
			option := languageAndOptions[i]
			optionParts := strings.Split(option, "=")

			if strings.ToUpper(optionParts[0]) == "Q" && len(optionParts) == 2 {
				var err error
				quality, err = strconv.ParseFloat(optionParts[1], 64)
				if err != nil {
					quality = 1
				}
			}
		}

		supported := false
		for _, supportedLanguage := range config.DrunsConfig.Languages {
			languageParts := strings.Split(language, "-")

			if strings.ToLower(language) == strings.ToLower(supportedLanguage) ||
				strings.ToLower(languageParts[0]) == strings.ToLower(supportedLanguage) {

				language = supportedLanguage
				supported = true
				break
			}
		}

		if supported && selectedQuality < quality {
			selectedLanguage = language
			selectedQuality = quality
		}
	}

	if selectedLanguage == "" && len(config.DrunsConfig.Languages) > 0 {
		selectedLanguage = config.DrunsConfig.Languages[0]
	}

	i.handler.SetLanguage(selectedLanguage)
	i.handler.SetMessages(tr.NewMessageHolder(selectedLanguage))
}
