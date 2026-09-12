package otp

type OTP interface {
	GenerateHash(username string) (string, string, error)
	ValidateCode(passcode, secretKey string) (bool, error)
}
