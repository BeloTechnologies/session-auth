package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/BeloTechnologies/session-core/core_models/user_models"
	"github.com/spf13/viper"
	"net/http"
	"session-auth/utils"
)

func CreateUserEntry(user user_models.CreateUserRow) (user_models.User, error) {
	log := utils.InitLogger() // Initialize and get the logger
	var successResponse core_models.SuccessResponse
	var userResponse user_models.User
	url := viper.GetString("proxies.user.url")

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Errorf("error marshalling user: %v", err)
		return userResponse, fmt.Errorf("error marshalling user: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/create_row/", url), bytes.NewBuffer(userJson))
	if err != nil {
		log.Errorf("error creating new request: %v", err)
		return userResponse, fmt.Errorf("error creating new request: %w", err)
	}

	// Set request headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error making POST request: %v", err)
		return userResponse, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Response validation
	if resp.StatusCode != http.StatusCreated {
		log.Errorf("unexpected status code: %d", resp.StatusCode)
		return userResponse, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		log.Errorf("error decoding response body: %v", err)
		return userResponse, fmt.Errorf("error decoding response body: %w", err)
	}

	// Unmarshal the Data field
	dataBytes, err := json.Marshal(successResponse.Data)
	if err != nil {
		log.Error("error marshalling proxy responses")
		return userResponse, fmt.Errorf("error marshalling proxy responses")
	}

	if err := json.Unmarshal(dataBytes, &userResponse); err != nil {
		log.Errorf("error unmarshalling user response: %v", err)
		return userResponse, fmt.Errorf("error unmarshalling user response: %w", err)
	}

	return userResponse, nil
}

// GetUser retrieves a user from the user service
func GetUser(ID int, token string) (user_models.User, error) {
	log := utils.InitLogger() // Initialize and get the logger
	var successResponse core_models.SuccessResponse
	var user user_models.User
	url := viper.GetString("proxies.user.url")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%d", url, ID), nil)
	if err != nil {
		log.Errorf("error creating new request: %v", err)
		return user, fmt.Errorf("error creating new request: %w", err)
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error making GET request: %v", err)
		return user, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Response validation
	if resp.StatusCode != http.StatusOK {
		log.Errorf("unexpected status code: %d", resp.StatusCode)
		return user, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&successResponse); err != nil {
		log.Errorf("error decoding response body: %v", err)
		return user, fmt.Errorf("error decoding response body: %w", err)
	}

	// Unmarshal the Data field into User
	dataBytes, err := json.Marshal(successResponse.Data)
	if err != nil {
		log.Error("error marshalling proxy responses")
		return user, fmt.Errorf("error marshalling proxy responses")
	}

	if err := json.Unmarshal(dataBytes, &user); err != nil {
		log.Errorf("error decoding response body: %v", err)
		return user, fmt.Errorf("error decoding response body: %w", err)
	}

	return user, nil
}
