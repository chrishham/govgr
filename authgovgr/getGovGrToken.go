package authgovgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetGovGrToken retrieves the token from the government API using the provided GSIS token and subdomain.
func getGovGrToken(gsisToken, govSubdomain string) (string, error) {
	apiURL := fmt.Sprintf("https://%s.services.gov.gr/api/token/?code=%s", govSubdomain, gsisToken)

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 35 * time.Second, // Timeout for the entire request
	}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	// Set request headers
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "el-GR,el;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Referer", fmt.Sprintf("https://%s.services.gov.gr/login/token/", govSubdomain))
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("TE", "trailers")

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	token, ok := result["token"].(string)
	if !ok || token == "" {
		return "", errors.New("token not found in response")
	}

	return token, nil
}
