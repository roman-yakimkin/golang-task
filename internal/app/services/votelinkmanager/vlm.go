package votelinkmanager

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"task/internal/app/errors"
	"task/internal/app/interfaces"
	"task/internal/app/services/configmanager"
)

type EncryptVoteLinkManager struct {
	config *configmanager.Config
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func NewEncryptVoteLinkManager(config *configmanager.Config) *EncryptVoteLinkManager {
	return &EncryptVoteLinkManager{
		config: config,
	}
}

func (vl *EncryptVoteLinkManager) encrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	encrypted := base64.StdEncoding.EncodeToString(cipherText)
	return hex.EncodeToString([]byte(encrypted)), nil
}

func (vl *EncryptVoteLinkManager) decrypt(text, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	decodedBytes, err := hex.DecodeString(text)
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(string(decodedBytes))
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func (vl *EncryptVoteLinkManager) Generate(data interfaces.VoteLinkData) string {
	str := fmt.Sprintf("%s %d %t %s", data.Email, data.TaskID, data.Result, data.Checksum)
	encLink, err := vl.encrypt(str, vl.config.SecretString)
	if err != nil {
		log.Fatal(err)
	}
	return encLink
}

func (vl *EncryptVoteLinkManager) Parse(link string) (interfaces.VoteLinkData, error) {
	decTex, err := vl.decrypt(link, vl.config.SecretString)
	if err != nil {
		return interfaces.VoteLinkData{}, err
	}
	var data interfaces.VoteLinkData
	_, err = fmt.Sscanf(decTex, "%s %d %t %s", &data.Email, &data.TaskID, &data.Result, &data.Checksum)
	if err != nil {
		return interfaces.VoteLinkData{}, errors.ErrVoteLinkNotParsed
	}
	return data, nil
}
