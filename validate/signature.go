package validate

import (
	"errors"
	"fmt"
	// "io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
)

func VerifySignature(filePath, sigPath, publicKeyPath string) error {
	keyFile, err := os.Open(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key: %w", err)
	}
	defer keyFile.Close()

	keyring, err := openpgp.ReadArmoredKeyRing(keyFile)
	if err != nil {
		return fmt.Errorf("failed to read public key: %w", err)
	}

	dataFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file to verify: %w", err)
	}
	defer dataFile.Close()

	
	sigFile, err := os.Open(sigPath)
	if err != nil {
		return fmt.Errorf("failed to open signature file: %w", err)
	}
	defer sigFile.Close()

	if _, err := openpgp.CheckArmoredDetachedSignature(keyring, dataFile, sigFile, nil); err != nil {
    return errors.New("signature verification failed: " + err.Error())
}


	return nil
}
