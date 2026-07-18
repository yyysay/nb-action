package pwd

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

const wireguardKeySize = 32

func generateWGKeypair() (map[string]interface{}, error) {

	privateKey, err := randomBytes(wireguardKeySize)

	if err != nil {
		return nil, fmt.Errorf(
			"generate private key: %w",
			err,
		)
	}

	// WireGuard Curve25519 clamp
	privateKey[0] &= 248
	privateKey[31] = (privateKey[31] & 127) | 64

	publicKey, err := curve25519.X25519(
		privateKey,
		curve25519.Basepoint,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"generate public key: %w",
			err,
		)
	}

	return map[string]interface{}{
		"private_key": base64.StdEncoding.EncodeToString(privateKey),
		"public_key":  base64.StdEncoding.EncodeToString(publicKey),
	}, nil
}
