package pwd

import "crypto/rand"

func randomBytes(size int) ([]byte, error) {

	buf := make([]byte, size)

	_, err := rand.Read(buf)

	return buf, err
}
