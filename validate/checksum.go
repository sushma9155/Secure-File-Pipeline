package validate

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strings"
)


func ValidateChecksum(filePath, expectedChecksum, algo string) error {
	hasher, err := getHasher(algo)
	if err != nil {
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to compute hash: %w", err)
	}

	computed := hex.EncodeToString(hasher.Sum(nil))
	expected := strings.ToLower(strings.TrimSpace(expectedChecksum))

	if computed != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, computed)
	}

	return nil
}

func getHasher(algo string) (hash.Hash, error) {
	switch strings.ToLower(algo) {
	case "md5":
		return md5.New(), nil
	case "sha256":
		return sha256.New(), nil
	default:
		return nil, fmt.Errorf("unsupported checksum algorithm: %s", algo)
	}
}
