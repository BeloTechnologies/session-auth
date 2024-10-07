package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"net/http"
	"net/http/httptest"
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

		assert.Equal(t, http.StatusCreated, w.Code)
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
		mockUser := bson.D{
			{Key: "email", Value: "test.user@example.com"},
			{Key: "password", Value: "hashedPassword"},
		}

		// Add the mock response for FindOne
		firstBatch := mtest.CreateCursorResponse(1, "sessionAuth.users", mtest.FirstBatch, mockUser)
		mt.AddMockResponses(firstBatch)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
