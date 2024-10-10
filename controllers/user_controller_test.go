package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("create user", func(mt *mtest.T) {
		r := gin.Default()

		r.POST("/users/create/", CreateUser(mt.Client))

		payload := `{
			"username": "testuser",
			"password": "password",	
			"email": "test.user@example.com",
			"phone": "1234567890"
		}`
		req, _ := http.NewRequest("POST", "/users/create/", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		r.ServeHTTP(w, req)

		checkSuccessResponse(t, w, http.StatusCreated)
	})
}

func TestLoginUser(t *testing.T) {
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
			{Key: "email", Value: "test.user@example.com"},
			{Key: "password", Value: hashed},
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
	var responseBody models.SuccessResponse
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
