package auth;
import (
	"net/http"
	"fmt"
	"strings"
)
func GetBearerToken(headers *http.Header) (string, error) {
	authValue := headers.Get("Authorization")
	if authValue == "" {
		return "", fmt.Errorf("No auth specified")
	}
	valueParts := strings.Fields(authValue)
	if len(valueParts) != 2 || valueParts[0] != "Bearer" {
		return "", fmt.Errorf("Malformed Authorization Header")
	}
	return valueParts[1], nil
}

func GetAPIKey(headers *http.Header) (string, error) {
        authValue := headers.Get("Authorization")
        if authValue == "" {
                return "", fmt.Errorf("No auth specified")
        }
        valueParts := strings.Fields(authValue)
        if len(valueParts) != 2 || valueParts[0] != "ApiKey" {
                return "", fmt.Errorf("Malformed Authorization Header")
        }
        return valueParts[1], nil
}
