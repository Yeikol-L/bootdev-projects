package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

type EncryptedPayload struct {
	Nonce      []byte
	AuthTag    []byte
	Ciphertext []byte
}

func encrypt(plaintext []byte, key [32]byte) (EncryptedPayload, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return EncryptedPayload{}, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptedPayload{}, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return EncryptedPayload{}, err
	}
	sealed := aesgcm.Seal(nil, nonce, plaintext, nil)
	tagSize := aesgcm.Overhead()
	ciphertext := sealed[:len(sealed)-tagSize]
	authTag := sealed[len(sealed)-tagSize:]
	return EncryptedPayload{
		Nonce:      nonce,
		AuthTag:    authTag,
		Ciphertext: ciphertext,
	}, nil

}

func decrypt(payload EncryptedPayload, key [32]byte) ([]byte, error) {
	if len(payload.Nonce) != 12 || len(payload.AuthTag) != 16 {
		return nil, errors.New("Bad options")
	}
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	cipher := []byte{}
	cipher = append(cipher, payload.Ciphertext...)
	cipher = append(cipher, payload.AuthTag...)
	return aesgcm.Open(nil, payload.Nonce, cipher, nil)
}
