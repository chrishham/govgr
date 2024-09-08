package authgovgr

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetGsisToken(t *testing.T) {

	const govSubdomain = "dilosi"
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	// Access environment variables
	gsisUserUsername := os.Getenv("gsisUserUsername")
	gsisUserPassword := os.Getenv("gsisUserPassword")

	gsisToken, err := getGsisToken(gsisUserUsername, gsisUserPassword, govSubdomain)
	if err != nil {
		t.Fatalf("Error getting GSIS token: %v", err)
	}

	if gsisToken == "" {
		t.Fatalf("Expected non-empty GSIS token, got an empty string")
	}

	fmt.Printf("GSIS token: %s\n", gsisToken)

	govGrToken, err := getGovGrToken(gsisToken, govSubdomain)

	if err != nil {
		t.Fatalf("Error getting GOV.GR token: %v", err)
	}

	if govGrToken == "" {
		t.Fatalf("Expected non-empty GOV.GR token, got an empty string")
	}

	fmt.Printf("GOV.GR token: %s\n", govGrToken)

	govGrToken1, err := GetGovGrTokenFromPool(gsisUserUsername, gsisUserPassword, govSubdomain)

	if err != nil {
		t.Fatalf("Error getting GOV.GR token 2: %v", err)
	}

	fmt.Printf("GOV.GR token: %s\n", govGrToken1)

	govGrToken2, err := GetGovGrTokenFromPool(gsisUserUsername, gsisUserPassword, govSubdomain)

	if err != nil {
		t.Fatalf("Error getting GOV.GR token 3: %v", err)
	}

	fmt.Printf("GOV.GR token: %s\n", govGrToken2)

	if govGrToken1 != govGrToken2 {
		t.Fatalf("GOV.GR token mismatch!")
	}
	t.Logf("Test passed")
}
