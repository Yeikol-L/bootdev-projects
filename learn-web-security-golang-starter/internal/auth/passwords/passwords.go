package passwords

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
)

const MaxLength = 128
const (
	memoryKiB     = 19 * 1024
	iterations    = 2
	parallelism   = 1
	derivedkeyLen = 32
)

func Hash(password string) (string, error) {
	if utf8.RuneCountInString(password) > MaxLength {
		return "", fmt.Errorf("password must not exceed %d characters", MaxLength)
	}
	salt, err := genSalt()
	if err != nil {
		return "", err
	}
	derivedKey := argon2.IDKey([]byte(password), salt[:], iterations, memoryKiB, parallelism, derivedkeyLen)
	hash := argon2idHash{
		version:     argon2.Version,
		memoryKiB:   memoryKiB,
		iterations:  iterations,
		parallelism: parallelism,
		salt:        salt[:],
		derivedKey:  derivedKey,
	}
	return encodeArgon2idHash(hash), nil
}

func Verify(password, encodedHash string) bool {
	if utf8.RuneCountInString(password) > MaxLength {
		return false
	}
	if !strings.HasPrefix(encodedHash, "$argon2id$") {
		return verifySha256(password, encodedHash)
	}
	hash, ok := parseArgon2idHash(encodedHash)
	if !ok || hash.version != argon2.Version {
		return false
	}
	candiate := argon2.IDKey([]byte(password), hash.salt, hash.iterations, hash.memoryKiB, hash.parallelism, derivedkeyLen)
	return subtle.ConstantTimeCompare(hash.derivedKey, candiate) == 1

}

func NeedsRehash(encodedHash string) bool {
	if !strings.HasPrefix(encodedHash, "$argon2id$") {
		return true
	}
	hash, ok := parseArgon2idHash(encodedHash)
	if !ok {
		return false
	}
	return hash.version != argon2.Version || hash.memoryKiB != memoryKiB || hash.iterations != iterations || hash.parallelism != parallelism || len(hash.derivedKey) != derivedkeyLen
}

func verifySha256(password, encodedHash string) bool {
	if utf8.RuneCountInString(password) > MaxLength {
		return false
	}
	expectedHash, ok := decodeLegacyHash(encodedHash)
	if !ok {
		return false
	}
	candidateHash := sha256.Sum256([]byte(password))
	return subtle.ConstantTimeCompare(candidateHash[:], expectedHash) == 1
}

func genSalt() ([16]byte, error) {
	salt := [16]byte{}
	if _, err := rand.Read(salt[:]); err != nil {
		return [16]byte{}, err
	}
	return salt, nil
}
