package pwd

import (
	"encoding/base64"

	"golang.org/x/crypto/curve25519"
)

func generateWGKeypair() (map[string]interface{}, error) {

	privateKey, err := randomBytes(32)

	if err != nil {
		return nil, err
	}

	// WireGuard clamp
	privateKey[0] &= 248
	privateKey[31] = (privateKey[31] & 127) | 64

	publicKey, err := curve25519.X25519(
		privateKey,
		curve25519.Basepoint,
	)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"private_key": base64.StdEncoding.EncodeToString(privateKey),
		"public_key":  base64.StdEncoding.EncodeToString(publicKey),
	}, nil
}
