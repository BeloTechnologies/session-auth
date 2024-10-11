package services

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"net/http"
	"session-auth/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("create user", func(mt *mtest.T) {
		user := models.CreateUser{
			Username: "testuser",
			Password: "password",
			Email:    "test.user@example.com",
			Phone:    "1234567890",
		}

		// Mock the document that should be returned by FindOne
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "sessionAuth.users", mtest.FirstBatch))
		// Mock the document that should be returned by InsertOne
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		tokenResponse, err := CreateUser(mt.Client, &user)

		assert.Nil(t, err)
		assert.NotNil(t, tokenResponse)
		assert.NotNil(t, tokenResponse.Token)
	})
}

func TestCreateUserUserExists(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("create user", func(mt *mtest.T) {
		user := models.CreateUser{
			Username: "testuser",
			Password: "password",
			Email:    "test.user@example.com",
			Phone:    "1234567890",
		}

		mockUser := bson.D{
			{Key: "email", Value: "test.user@example.com"},
		}

		// Mock the document that should be returned by FindOne
		firstBatch := mtest.CreateCursorResponse(1, "sessionAuth.users", mtest.FirstBatch, mockUser)
		mt.AddMockResponses(firstBatch)

		tokenResponse, err := CreateUser(mt.Client, &user)

		assert.Nil(t, tokenResponse)
		assert.NotNil(t, err)
		assert.Equal(t, http.StatusConflict, err.Status)
	})
}

func TestLoginUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("login user", func(mt *mtest.T) {
		user := models.LoginUser{
			Email:    "test.user@example.com",
			Password: "password",
		}

		// Mock the document that should be returned by FindOne
		hashed, err := HashPassword("password")
		assert.NoError(t, err)

		mockUser := bson.D{
			{Key: "email", Value: "test.user@example.com"},
			{Key: "password", Value: hashed},
		}

		// Add the mock response for FindOne
		firstBatch := mtest.CreateCursorResponse(1, "sessionAuth.users", mtest.FirstBatch, mockUser)
		mt.AddMockResponses(firstBatch)

		tokenResponse, sessionError := LoginUser(mt.Client, &user)

		assert.Nil(t, sessionError)
		assert.NotNil(t, tokenResponse)
		assert.NotNil(t, tokenResponse.Token)
	})
}

func TestLoginUserNotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("login user", func(mt *mtest.T) {
		user := models.LoginUser{
			Email:    "test.user@example.com",
			Password: "password",
		}

		// Add the mock response for FindOne
		firstBatch := mtest.CreateCursorResponse(0, "sessionAuth.users", mtest.FirstBatch)
		mt.AddMockResponses(firstBatch)

		tokenResponse, sessionError := LoginUser(mt.Client, &user)

		assert.Nil(t, tokenResponse)
		assert.NotNil(t, sessionError)
		assert.Equal(t, http.StatusNotFound, sessionError.Status)
	})
}

func TestLoginUserInvalidPassword(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("login user", func(mt *mtest.T) {
		user := models.LoginUser{
			Email:    "test.user@example.com",
			Password: "invalid_password",
		}

		// Mock the document that should be returned by FindOne
		hashed, err := HashPassword("password")
		assert.NoError(t, err)

		mockUser := bson.D{
			{Key: "email", Value: "test.user@example.com"},
			{Key: "password", Value: hashed},
		}

		// Add the mock response for FindOne
		firstBatch := mtest.CreateCursorResponse(1, "sessionAuth.users", mtest.FirstBatch, mockUser)
		mt.AddMockResponses(firstBatch)

		tokenResponse, sessionError := LoginUser(mt.Client, &user)

		assert.Nil(t, tokenResponse)
		assert.NotNil(t, sessionError)
		assert.Equal(t, http.StatusBadRequest, sessionError.Status)
	})
}
