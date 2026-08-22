package hashing

type Hashing interface {
	Encrypt(plainText string) (string, error)
	Decrypt(secureText string) (string, error)
}
