package validation

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func ValidateHash(filePath string, expectedHash string) (bool, error) {
	hash, err := CalculateHash(filePath)
	if err != nil {
		return false, err
	}
	return hash == expectedHash, nil
}

