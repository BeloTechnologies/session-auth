package services

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"session-auth/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("create user", func(mt *mtest.T) {
		user := models.User{
			Username: "testuser",
			Password: "password",
			Email:    "test.user@example.com",
			Phone:    "1234567890",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())

		userID, err := CreateUser(mt.Client, &user)

		assert.NoError(t, err)
		assert.NotNil(t, userID)
	})
}
