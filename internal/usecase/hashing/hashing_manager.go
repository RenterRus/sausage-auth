package hashing

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type hashManager struct {
	key []byte
}

func NewHashingManager(key []byte) Hashing {
	return &hashManager{
		key: key,
	}
}

// Encrypt шифрует строку и возвращает Base64-код
func (h *hashManager) Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher(h.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherBytes := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// Decrypt расшифровывает Base64-строку обратно в текст
func (h *hashManager) Decrypt(secureText string) (string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(secureText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(h.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", fmt.Errorf("ошибка: текст слишком короткий")
	}

	nonce, cipherText := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plainBytes, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainBytes), nil
}
