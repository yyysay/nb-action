package bark

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	t.Setenv(
		"BARK_AES_KEY",
		"12345678901234567890123456789012",
	)

	t.Setenv(
		"BARK_AES_IV",
		"123456789012",
	)

	t.Setenv(
		"BARK_SERVER",
		"https://api.day.app",
	)

	t.Setenv(
		"BARK_DEVICE_KEY",
		"test-device-key",
	)

	config, err := LoadConfig()

	if err != nil {

		t.Fatal(err)

	}

	if len(config.Key) != 32 {

		t.Fatalf(
			"key length error: %d",
			len(config.Key),
		)

	}

	if len(config.IV) != 12 {

		t.Fatalf(
			"iv length error: %d",
			len(config.IV),
		)

	}

	if config.Server == "" {

		t.Fatal(
			"server empty",
		)

	}

	if config.DeviceKey == "" {

		t.Fatal(
			"device key empty",
		)

	}

	t.Log(
		"config load success",
	)

}
