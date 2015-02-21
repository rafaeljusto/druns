package interceptor

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/rafaeljusto/druns/Godeps/_workspace/src/github.com/gustavo-hms/trama"
	"github.com/rafaeljusto/druns/core/tr"
	"github.com/rafaeljusto/druns/web/config"
)

type languageMessager interface {
	SetLanguage(language string)
	SetMessages(messageHolder tr.MessageHolder)
}

////////////////////////////////////////////////////////////
/////////////////////// AJAX ///////////////////////////////
////////////////////////////////////////////////////////////

type AcceptLanguageAJAX struct {
	trama.NopAJAXInterceptor
	handler languageMessager
}

func NewAcceptLanguageAJAX(h languageMessager) *AcceptLanguageAJAX {
	return &AcceptLanguageAJAX{handler: h}
}

func (i AcceptLanguageAJAX) Before(w http.ResponseWriter, r *http.Request) {
	selectedLanguage := acceptLanguage(r)
	i.handler.SetLanguage(selectedLanguage)
	i.handler.SetMessages(tr.NewMessageHolder(selectedLanguage))
}

////////////////////////////////////////////////////////////
/////////////////////// WEB ////////////////////////////////
////////////////////////////////////////////////////////////

type AcceptLanguageWeb struct {
	trama.NopWebInterceptor
	handler languageMessager
}

func NewAcceptLanguageWeb(h languageMessager) *AcceptLanguageWeb {
	return &AcceptLanguageWeb{handler: h}
}

func (i AcceptLanguageWeb) Before(response trama.Response, r *http.Request) {
	selectedLanguage := acceptLanguage(r)
	response.SetTemplateGroup(selectedLanguage)
	i.handler.SetLanguage(selectedLanguage)
	i.handler.SetMessages(tr.NewMessageHolder(selectedLanguage))
}

////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////

func acceptLanguage(r *http.Request) string {
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

	return selectedLanguage
}
