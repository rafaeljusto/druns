package password

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/rafaeljusto/druns/core"
)

func Encrypt(password string) (string, error) {
	block, err := aes.NewCipher(getKey())
	if err != nil {
		return "", core.NewError(err)
	}

	iv := make([]byte, block.BlockSize())
	if _, err = rand.Read(iv); err != nil {
		return "", core.NewError(err)
	}

	ofbStream := cipher.NewOFB(block, iv)
	output := make([]byte, len(password))
	ofbStream.XORKeyStream(output, []byte(password))

	buffer := bytes.NewBuffer(iv)
	buffer.Write(output)
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

func Decrypt(password string) (string, error) {
	encryptedPassword, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		return "", core.NewError(err)
	}

	block, err := aes.NewCipher(getKey())
	if err != nil {
		return "", core.NewError(err)
	}

	if len(encryptedPassword) < block.BlockSize() {
		return "", core.NewError(errors.New("password is to small to be decrypted"))
	}

	iv := encryptedPassword[:block.BlockSize()]
	encryptedPassword = encryptedPassword[block.BlockSize():]

	ofbStream := cipher.NewOFB(block, iv)
	output := make([]byte, len(encryptedPassword))
	ofbStream.XORKeyStream(output, encryptedPassword)
	return string(output), nil
}
