package protocol

import (
	"strings"
)

const (
	MsgCodeCorruptedData = "corrupted-data"
)

type msgCode string

type Translator interface {
	Translate(lang string) bool
}

type Message struct {
	Code  msgCode `json:"code"`
	Field string  `json:"field,omitempty"`
	Value string  `json:"value,omitempty"`
	Text  string  `json:"message,omitempty"`
}

func NewMessage(code msgCode, value string) Translator {
	return &Message{
		Code:  code,
		Value: value,
	}
}

func NewMessageWithField(code msgCode, field string, value string) Translator {
	return &Message{
		Code:  code,
		Field: field,
		Value: value,
	}
}

func (msg *Message) Translate(language string) bool {
	if len(language) == 0 || language == "*" {
		language = "en"

	} else {
		language = strings.Split(language, "-")[0]
	}

	if translation, ok := Translations[language]; ok {
		msg.Text = translation[msg.Code]
		return true
	}

	return false
}

type MessagesHolder struct {
	msgs Messages
}

func (m *MessagesHolder) Add(trs ...Translator) {
	for _, tr := range trs {
		if tr == nil {
			return
		}

		switch tr.(type) {
		case *Message:
			m.msgs = append(m.msgs, tr.(*Message))
		case Messages:
			m.msgs = append(m.msgs, tr.(Messages)...)
		}
	}
}

func (m *MessagesHolder) Messages() Translator {
	if len(m.msgs) == 0 {
		return nil
	}

	return m.msgs
}

type Messages []Translator

func (m Messages) Translate(language string) bool {
	for i := range m {
		if ok := m[i].Translate(language); !ok {
			return false
		}
	}

	return true
}
