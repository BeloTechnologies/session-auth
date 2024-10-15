package proxy

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

func CreateUserEntryInUserProxy() (*http.Response, error) {
	url := viper.GetString("proxies.user.url")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/create/", url), bytes.NewBuffer(nil))
	if err != nil {
		return nil, fmt.Errorf("error creating new request: %w", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}

	return resp, nil
}
