package utils

import (
	"crypto/aes"
	"encoding/hex"
)

func EncryptAES(key []byte, plaintext string) (string, error) {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	for i := 1; i <= len(plaintext)/16; i++ {
		tempBuf := make([]byte, 16)

		offset := (i - 1) * 16
		limit := offset + 16
		// encrypt
		c.Encrypt(tempBuf, []byte(plaintext[offset:limit]))

		for j := 0; j < len(tempBuf); j++ {
			out[offset+j] = tempBuf[j]
		}
	}
	// return hex string
	return hex.EncodeToString(out), nil
}

func DecryptAES(key []byte, ct string) (string, error) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s, nil
}
