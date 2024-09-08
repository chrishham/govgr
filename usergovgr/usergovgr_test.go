package usergovgr

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetUserInfo(t *testing.T) {

	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	// Access environment variables
	gsisUserUsername := os.Getenv("gsisUserUsername")
	gsisUserPassword := os.Getenv("gsisUserPassword")

	// Call the actual getUserInfo function
	userInfo, err := getUserInfo(gsisUserUsername, gsisUserPassword)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the returned userInfo
	if userInfo == nil {
		t.Fatalf("Expected userInfo, got nil")
	}

	if userInfo.FirstName == "" || userInfo.Surname == "" || userInfo.AFM == "" {
		t.Errorf("Expected valid user info, got: %+v", userInfo)
	}

	// Print the user info
	fmt.Printf("User Info:\n")
	fmt.Printf("First Name: %s\n", userInfo.FirstName)
	fmt.Printf("Surname: %s\n", userInfo.Surname)
	fmt.Printf("AFM: %s\n", userInfo.AFM)
	fmt.Printf("Birth Date: %s\n", userInfo.BirthDate)
	fmt.Printf("Mobile Certified Login: %s\n", userInfo.MobileCertifiedLogin)
}
