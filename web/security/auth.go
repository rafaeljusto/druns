package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/rafaeljusto/druns/core/protocol"
)

const (
	namespace = "druns"
)

func AddAuthorization(r *http.Request, handle, secretId, secret string) protocol.Translator {
	stringToSign, err := buildStringToSign(r, secretId)
	if err != nil {
		return err
	}

	signature := generateSignature(stringToSign, secret)
	r.Header.Set("Authorization", fmt.Sprintf("%s %s:%s:%s", namespace, handle, secretId, signature))
	return nil
}

func CheckAuthorization(r *http.Request, secret func(string) (string, error)) (bool, string, protocol.Translator, error) {
	authorization := r.Header.Get("Authorization")
	authorization = strings.TrimSpace(authorization)

	if len(authorization) == 0 {
		return false, "", protocol.NewMessageWithField(protocol.MsgCodeMissingData,
			"Authorization", ""), nil
	}

	// Authorization format: "<namespace> <handle>:<secretId>:<secret>"
	authorizationParts := strings.Split(authorization, " ")
	if len(authorizationParts) != 2 {
		return false, "", protocol.NewMessageWithField(protocol.MsgCodeInvalidFormat,
			"Authorization", authorization), nil
	}

	ns := authorizationParts[0]
	ns = strings.TrimSpace(ns)
	ns = strings.ToLower(ns)

	if ns != namespace {
		return false, "", protocol.NewMessageWithField(protocol.MsgCodeInvalidFormat,
			"Authorization", authorization), nil
	}

	secretParts := strings.Split(authorizationParts[1], ":")
	if len(secretParts) != 3 {
		return false, "", protocol.NewMessageWithField(protocol.MsgCodeInvalidFormat,
			"Authorization", authorization), nil
	}

	handle := secretParts[0]
	secretId := secretParts[1]
	secretId = strings.TrimSpace(secretId)
	secretId = strings.ToLower(secretId)

	stringToSign, msg := buildStringToSign(r, secretId)
	if msg != nil {
		return false, handle, msg, nil
	}

	s, err := secret(secretId)
	if err != nil {
		return false, handle, nil, err
	}

	signature := generateSignature(stringToSign, s)
	if signature != secretParts[2] {
		return false, handle, nil, nil
	}

	return true, handle, nil, nil
}

func buildStringToSign(r *http.Request, secretId string) (string, protocol.Translator) {
	// StringToSign = HTTP-Verb + "\n" +
	//  Content-MD5 + "\n" + // RFC1864
	//  Content-Type + "\n" +
	//  Date + "\n" +
	//  AccessKeyID + "\n" +
	//  Path + "\n" +
	//  CanonicalizedQueryString;

	stringToSign := r.Method

	contentMD5 := ""
	contentType := ""

	if r.ContentLength > 0 {
		contentMD5 = getHTTPContentMD5(r)

		if len(contentMD5) == 0 {
			return "", protocol.NewMessageWithField(protocol.MsgCodeMissingData, "Content-Md5", "")
		}

		contentType = getHTTPContentType(r)

		if len(contentType) == 0 {
			return "", protocol.NewMessageWithField(protocol.MsgCodeMissingData, "Content-Type", "")
		}

		// For now we are ignoring version
		if idx := strings.Index(contentType, ";"); idx > 0 {
			contentType = contentType[0:idx]
		}

		stringToSign = fmt.Sprintf("%s\n%s", stringToSign, contentMD5)
		stringToSign = fmt.Sprintf("%s\n%s", stringToSign, contentType)
	}

	dateStr := getHTTPDate(r)
	if len(dateStr) == 0 {
		return "", protocol.NewMessageWithField(protocol.MsgCodeMissingData, "Date", "")
	}

	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, dateStr)
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, secretId)
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, r.URL.Path)

	var queryString []string
	for key, values := range r.URL.Query() {
		for _, value := range values {
			keyAndValue := fmt.Sprintf("%s=%s", key, value)
			queryString = append(queryString, keyAndValue)
		}
	}

	sort.Strings(queryString)
	sortedQueryString := strings.Join(queryString, "&")
	stringToSign = fmt.Sprintf("%s\n%s", stringToSign, sortedQueryString)

	return stringToSign, nil
}

func generateSignature(stringToSign, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getHTTPContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	contentType = strings.TrimSpace(contentType)
	contentType = strings.ToLower(contentType)
	return contentType
}

func getHTTPContentMD5(r *http.Request) string {
	contentMD5 := r.Header.Get("Content-MD5")
	contentMD5 = strings.TrimSpace(contentMD5)
	return contentMD5
}

func getHTTPDate(r *http.Request) string {
	dateStr := r.Header.Get("Date")
	dateStr = strings.TrimSpace(dateStr)
	return dateStr
}
