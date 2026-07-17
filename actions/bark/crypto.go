package bark

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func encrypt(
	payload BarkPayload,
) (
	string,
	error,
) {

	key := os.Getenv("BARK_AES_KEY")
	iv := os.Getenv("BARK_AES_IV")

	if len(key) != 32 || len(iv) != 12 {
		return "",
			fmt.Errorf(
				"BARK_AES_KEY(32位) 或 BARK_AES_IV(12位) 环境变量配置错误",
			)
	}

	data, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(
		[]byte(key),
	)

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(
		nil,
		[]byte(iv),
		data,
		nil,
	)

	return url.QueryEscape(
		base64.StdEncoding.EncodeToString(encrypted),
	), nil
}
