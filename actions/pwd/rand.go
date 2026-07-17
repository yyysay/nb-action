package pwd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
)

func generateRand(args []string) (map[string]interface{}, error) {

	if len(args) == 0 {
		return nil, fmt.Errorf("rand requires size")
	}

	size, err := strconv.Atoi(args[0])
	if err != nil || size <= 0 {
		return nil, fmt.Errorf("invalid rand size: %s", args[0])
	}

	data, err := randomBytes(size)
	if err != nil {
		return nil, fmt.Errorf("generate random data failed: %w", err)
	}

	value := hex.EncodeToString(data)

	for _, arg := range args[1:] {
		if arg == "--base64" {
			value = base64.StdEncoding.EncodeToString(data)
			break
		}
	}

	return map[string]interface{}{
		"value": value,
	}, nil
}
