package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/BeloTechnologies/session-core/core_models/user_models"
	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"log"
	"net/http"
	"net/http/httptest"
	"session-auth/models"
	"session-auth/services"
	"testing"
)

func TestCreateUser(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	user := user_models.CreateUserRowResponse{
		ID:        1,
		Username:  "john.doe",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@example.com",
		Phone:     "123-456-7890",
	}

	successResponse := core_models.SuccessResponse{
		Message: "User created successfully",
		Status:  http.StatusCreated,
		Data:    user,
	}

	successJson, err := json.Marshal(successResponse)
	assert.NoError(t, err)

	httpmock.RegisterResponder("POST", viper.GetString("proxies.user.url")+"/users/create_row/",
		httpmock.NewStringResponder(http.StatusCreated, string(successJson)))

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("create user", func(mt *mtest.T) {
		r := gin.Default()

		r.POST("/users/create/", CreateUser(mt.Client))

		inputtedUser := models.CreateUser{
			Username:  "john.doe",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
			Email:     "johndoe@example.com",
			Phone:     "123-456-7890",
		}

		userJson, err := json.Marshal(inputtedUser)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/users/create/", bytes.NewBuffer(userJson))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		// Mock the document that should be returned by FindOne
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "sessionAuth.users", mtest.FirstBatch))
		// Mock the document that should be returned by InsertOne
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		r.ServeHTTP(w, req)

		checkSuccessResponse(t, w, http.StatusCreated)
	})
}

func TestLoginUser(t *testing.T) {
	testID := 1
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	user := user_models.User{
		Username:       "john.doe",
		FirstName:      "John",
		LastName:       "Doe",
		Email:          "johndoe@example.com",
		Phone:          "123-456-7890",
		CreatedAt:      "2021-01-01T00:00:00Z",
		FollowersCount: 0,
		FollowingCount: 0,
	}

	successResponse := core_models.SuccessResponse{
		Message: "User created successfully",
		Status:  http.StatusOK,
		Data:    user,
	}

	successJson, err := json.Marshal(successResponse)
	assert.NoError(t, err)

	httpmock.RegisterResponder("GET", viper.GetString("proxies.user.url")+fmt.Sprintf("/users/%d", testID),
		httpmock.NewStringResponder(http.StatusOK, string(successJson)))

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("login user", func(mt *mtest.T) {
		r := gin.Default()

		r.POST("/users/login/", LoginUser(mt.Client))

		payload := `{
			"email": "test.user@example.com",
			"password": "password"
		}`
		req, _ := http.NewRequest("POST", "/users/login/", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		// Mock the document that should be returned by FindOne
		hashed, err := services.HashPassword("password")
		assert.NoError(t, err)

		mockUser := bson.D{
			{Key: "email", Value: "johndoe@example.com"},
			{Key: "password", Value: hashed},
			{Key: "username", Value: "john.doe"},
			{Key: "first_name", Value: "John"},
			{Key: "last_name", Value: "Doe"},
			{Key: "phone", Value: "123-456-7890"},
			{Key: "psql_id", Value: testID},
		}

		// Add the mock response for FindOne
		firstBatch := mtest.CreateCursorResponse(1, "sessionAuth.users", mtest.FirstBatch, mockUser)
		mt.AddMockResponses(firstBatch)

		r.ServeHTTP(w, req)

		checkSuccessResponse(t, w, http.StatusOK)
	})
}

func checkSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	log.Println("Checking user controller success responses")

	assert.Equal(t, expectedCode, w.Code)

	// Unmarshal the response body into the SuccessResponse struct
	var responseBody core_models.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.NoError(t, err)

	assert.Equal(t, expectedCode, responseBody.Status)

	// Unmarshal the Data field into AuthResponse
	dataBytes, err := json.Marshal(responseBody.Data)
	assert.NoError(t, err)

	var authResponse models.AuthResponse
	err = json.Unmarshal(dataBytes, &authResponse)
	assert.NoError(t, err)

	assert.NotEmpty(t, authResponse.Token)
}
