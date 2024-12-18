package auth

import (
	"time"
	"testing"
	"github.com/google/uuid"
	"github.com/madhu1992blue/go-servers-demo/internal/auth"
)

func TestJWT(t *testing.T) {
	origUUID := uuid.New()
	tokenSecret := "Thisisatokensecret"
	tokenString, err := auth.MakeJWT(origUUID, tokenSecret, 2 * time.Minute)
	if err != nil {
		t.Fatalf("Error Making JWT: %v", err)
	}
	gotUUID, err := auth.ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}
	if origUUID.String() != gotUUID.String() {
		t.Fatalf("UUID match - Got: %s, Want: %s", gotUUID,origUUID)
	}

	wrongTokenSecret := "Thisisawrongsecret"
	gotUUID, err = auth.ValidateJWT(tokenString, wrongTokenSecret)
        if err == nil {
		t.Fatalf("Error : JWT validation was successful with wrong secret")
        }
	if gotUUID != (uuid.UUID{}) {
		t.Fatalf("Error: Shouldn't get a valid UUID when wrong secret used")
	}

}
