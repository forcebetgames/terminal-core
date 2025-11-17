package domain

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetIPInfo(terminal *Terminal) error {
	// Make a request to the ipinfo.io API
	resp, err := http.Get("https://ipinfo.io/json")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Unmarshal the JSON response into the IPInfo struct
	if err := json.Unmarshal(body, &terminal); err != nil {
		return err
	}

	return nil
}
