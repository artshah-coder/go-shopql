package password

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

func HashPass(plainPassword string, salt []byte) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	res := make([]byte, len(salt))
	copy(res, salt)
	return append(res, hashedPass...)
}

func MakeSalt(n uint32) []byte {
	salt := make([]byte, n)
	rand.Read(salt)
	return salt
}
