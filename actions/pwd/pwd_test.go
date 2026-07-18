package pwd

import (
	"context"
	"testing"
)

func TestPassword_Name(t *testing.T) {

	p := New()

	if p.Name() != "pwd" {
		t.Fatalf(
			"expected pwd, got %s",
			p.Name(),
		)
	}
}

func TestPassword_Rand(t *testing.T) {

	p := New()

	result, err := p.Execute(
		context.Background(),
		[]string{
			"rand",
			"16",
		},
		nil,
	)

	if err != nil {
		t.Fatalf(
			"rand failed: %v",
			err,
		)
	}

	value, ok := result["value"].(string)

	if !ok {
		t.Fatal(
			"value missing",
		)
	}

	// hex编码，一个byte=2字符
	if len(value) != 32 {
		t.Fatalf(
			"expected length 32, got %d",
			len(value),
		)
	}
}

func TestPassword_RandBase64(t *testing.T) {

	p := New()

	result, err := p.Execute(
		context.Background(),
		[]string{
			"rand",
			"16",
			"--base64",
		},
		nil,
	)

	if err != nil {
		t.Fatalf(
			"rand base64 failed: %v",
			err,
		)
	}

	if result["value"] == "" {
		t.Fatal(
			"empty value",
		)
	}
}

func TestPassword_WireGuard(t *testing.T) {

	p := New()

	result, err := p.Execute(
		context.Background(),
		[]string{
			"wg-keypair",
		},
		nil,
	)

	if err != nil {
		t.Fatalf(
			"wireguard failed: %v",
			err,
		)
	}

	privateKey, ok := result["private_key"].(string)

	if !ok || privateKey == "" {
		t.Fatal(
			"missing private key",
		)
	}

	publicKey, ok := result["public_key"].(string)

	if !ok || publicKey == "" {
		t.Fatal(
			"missing public key",
		)
	}
}

func TestPassword_InvalidCommand(t *testing.T) {

	p := New()

	_, err := p.Execute(
		context.Background(),
		[]string{
			"unknown",
		},
		nil,
	)

	if err == nil {
		t.Fatal(
			"expected error",
		)
	}
}
