package authgovgr

import (
	"fmt"
	"sync"
	"time"
)

// TokenInfo holds information about a token and its expiration.
type TokenInfo struct {
	GsisUserUsername string
	GovSubdomain     string
	GovGrToken       string
	Expiration       time.Time
}

// Global pool to store tokens
var (
	tokenPool = make(map[string]TokenInfo)
	mu        sync.Mutex
)

// GetGovGrTokenFromPool retrieves a GovGr token from a pool. If the token is expired or doesn't exist, it retrieves a new one.
func GetGovGrTokenFromPool(gsisUserUsername, gsisUserPassword, govSubdomain string) (string, error) {
	key := fmt.Sprintf("%s:%s", gsisUserUsername, govSubdomain)
	mu.Lock()
	defer mu.Unlock()

	// Check if the token is in the pool and still valid
	if tokenInfo, exists := tokenPool[key]; exists && tokenInfo.Expiration.After(time.Now()) {
		return tokenInfo.GovGrToken, nil
	}

	fmt.Println("Request new GovGrToken")

	// Retrieve GSIS code
	gsisToken, err := getGsisToken(gsisUserUsername, gsisUserPassword, govSubdomain)
	if err != nil {
		return "", err
	}

	fmt.Println("New gsisToken retreived!")

	// Retrieve GovGr token
	govGrToken, err := getGovGrToken(gsisToken, govSubdomain)
	if err != nil {
		return "", err
	}
	fmt.Println("New govGrToken retreived!")

	// Update token pool with the new token and expiration time
	expiration := time.Now().Add(15 * time.Minute)

	tokenPool[key] = TokenInfo{
		GsisUserUsername: gsisUserUsername,
		GovSubdomain:     govSubdomain,
		GovGrToken:       govGrToken,
		Expiration:       expiration,
	}

	return govGrToken, nil
}
