package services

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
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

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		tokenResponse, err := CreateUser(mt.Client, &user)

		assert.NoError(t, err)
		assert.NotNil(t, tokenResponse)
		assert.NotNil(t, tokenResponse.Token)
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
		mockUser := bson.D{
			{Key: "email", Value: "test.user@example.com"},
			{Key: "password", Value: "hashedPassword"},
		}

		// Add the mock response for FindOne
		firstBatch := mtest.CreateCursorResponse(1, "sessionAuth.users", mtest.FirstBatch, mockUser)
		mt.AddMockResponses(firstBatch)

		tokenResponse, err := LoginUser(mt.Client, &user)

		assert.NoError(t, err)
		assert.NotNil(t, tokenResponse)
		assert.NotNil(t, tokenResponse.Token)
	})
}
