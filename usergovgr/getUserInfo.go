package usergovgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chrishham/govgr/authgovgr"
)

// Define the structure of the user info response
type UserInfo struct {
	MobileCertifiedLogin string `json:"mobile_certified_login"`
	FirstName            string `json:"firstname"`
	Surname              string `json:"surname"`
	AFM                  string `json:"afm"`
	BirthDate            string `json:"birth_date"`
}

const UA = "Mozilla/5.0 (Windows NT 10.0; rv:95.0) Gecko/20100101 Firefox/95.0"

func getUserInfo(gsisUserUsername, gsisUserPassword string) (*UserInfo, error) {

	govGrToken, err := authgovgr.GetGovGrTokenFromPool(gsisUserUsername, gsisUserPassword, "dilosi")
	if err != nil {
		return nil, err
	}

	url := "https://dilosi.services.gov.gr/api/users/me/"

	client := &http.Client{
		Timeout: 15 * time.Second, // Set a timeout for the request
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UA)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "el-GR,el;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", govGrToken))
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("TE", "trailers")

	// Create a context with timeout for the request
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Attach context to the request
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo UserInfo
	if err := json.Unmarshal(bodyBytes, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
