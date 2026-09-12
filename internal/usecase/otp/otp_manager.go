package otp

import (
	"fmt"

	"github.com/RenterRus/sausage-auth/internal/usecase/hashing"
	"github.com/pquerna/otp/totp"
)

type otpManager struct {
	// Ключ должен быть строго 16, 24 или 32 байта (для AES-128, 192 или 256)
	hash   hashing.Hashing
	issuer string
}

func NewOTPManager(hashing hashing.Hashing, issuer string) OTP {
	return &otpManager{
		hash:   hashing,
		issuer: issuer,
	}
}

func (o *otpManager) GenerateHash(username string) (string, string, error) {
	// Генерируем ключ
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      o.issuer,
		AccountName: username,
	})
	if err != nil {
		return "", "", err
	}

	// Secret() — строка для сохранения в БД (в зашифрованном виде!)
	secretKey, err := o.hash.Encrypt(key.Secret())
	if err != nil {
		return "", "", fmt.Errorf("GenerateHash.encrypt: %w", err)
	}

	// URL для построения QR-кода, который сканирует пользователь в приложении
	return secretKey, key.URL(), nil
}

func (o *otpManager) ValidateCode(passcode, secretKey string) (bool, error) {
	secret, err := o.hash.Decrypt(secretKey)
	if err != nil {
		return false, fmt.Errorf("ValidateCode.decrypt: %w", err)
	}

	// Валидирует код по текущему времени сервера
	return totp.Validate(passcode, secret), nil
}
