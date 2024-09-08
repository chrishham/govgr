package authgovgr

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; rv:95.0) Gecko/20100101 Firefox/95.0"

func getGsisToken(gsisUserUsername, gsisUserPassword, govSubdomain string) (string, error) {

	// Create a cookie jar to store cookies between requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}

	// Create an HTTP client with the cookie jar
	client := &http.Client{
		Jar:     jar,
		Timeout: 35 * time.Second,
	}

	// Step 1: Visit initial GSIS login screen to get cookies
	loginURL := fmt.Sprintf("https://%s.services.gov.gr/api/login/?backend=gsis", govSubdomain)
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "el-GR,el;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://auth.services.gov.gr/")
	req.Header.Set("TE", "trailers")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Step 2: Perform the login at GSIS
	loginURL = "https://oauth2.gsis.gr/oauth2server/j_spring_security_check"
	data := url.Values{}
	data.Set("j_username", gsisUserUsername)
	data.Set("j_password", gsisUserPassword)

	req, err = http.NewRequest("POST", loginURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "el-GR,el;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", "https://oauth2.gsis.gr")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://oauth2.gsis.gr/oauth2server/login.jsp")

	resp, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()

	if strings.Contains(finalURL, "/login/token/#") {
		return extractToken(finalURL), nil
	} else if strings.Contains(finalURL, "error_code=") {
		return "", errors.New(parseErrorCode(finalURL))
	} else if strings.Contains(finalURL, "authentication_error") {
		return "", errors.New("authentication_error")
	} else if strings.Contains(finalURL, "https://oauth2.gsis.gr/oauth2server/oauth/authorize") {
		// Step 3: Authorize gov.gr
		authorizeURL := "https://oauth2.gsis.gr/oauth2server/oauth/authorize"
		data = url.Values{}
		data.Set("user_oauth_approval", "true")
		data.Set("scope.read", "true")

		req, err = http.NewRequest("POST", authorizeURL, bytes.NewBufferString(data.Encode()))
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "el-GR,el;q=0.8,en-US;q=0.5,en;q=0.3")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Origin", "https://oauth2.gsis.gr")
		req.Header.Set("DNT", "1")
		req.Header.Set("Connection", "keep-alive")

		resp, err = client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		finalURL = resp.Request.URL.String()

		if strings.Contains(finalURL, "/login/token/#") {
			return extractToken(finalURL), nil
		} else {
			return "", errors.New("couldn't retrieve token, even after gov.gr authorization")
		}
	} else {
		return "", errors.New(finalURL)
	}
}

func extractToken(url string) string {
	// Find the substring after "token/#"
	prefix := "token/#"
	idx := strings.Index(url, prefix)

	// Extract everything after "token/#"
	token := url[idx+len(prefix):]
	return token
}

func parseErrorCode(url string) string {
	const errorCodePrefix = "error_code="
	index := len(url) - len(errorCodePrefix)
	return url[index:]
}
