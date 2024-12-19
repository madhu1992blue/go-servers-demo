package auth
import (
	"encoding/hex"
	"crypto/rand"
)
func MakeRefreshToken() (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}


