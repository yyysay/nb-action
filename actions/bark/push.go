package bark

import (
	"fmt"
	"net/http"
)

func push(
	config Config,
	ciphertext string,
) error {

	requestURL := fmt.Sprintf(
		"%s/%s/%s",
		config.Server,
		config.DeviceKey,
		ciphertext,
	)

	request, err := http.NewRequest(
		http.MethodGet,
		requestURL,
		nil,
	)

	if err != nil {

		return err

	}

	client := &http.Client{}

	response, err := client.Do(
		request,
	)

	if err != nil {

		return err

	}

	defer response.Body.Close()

	if response.StatusCode < 200 ||
		response.StatusCode >= 300 {

		return fmt.Errorf(
			"bark push failed: %s",
			response.Status,
		)

	}

	return nil
}
