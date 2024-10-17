package services

import (
	"encoding/json"
	"fmt"
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/BeloTechnologies/session-core/core_models/user_models"
	"github.com/jarcoal/httpmock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"net/http"
	"session-auth/models"
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
		inputtedUser := models.CreateUser{
			Username: "john.doe",
			Password: "password",
			Email:    "johndoe@example.com",
			Phone:    "123-456-7890",
		}

		// Mock the document that should be returned by FindOne
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "sessionAuth.users", mtest.FirstBatch))
		// Mock the document that should be returned by InsertOne
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		authResponse, err := CreateUser(mt.Client, &inputtedUser)

		assert.Nil(t, err)
		assert.NotNil(t, authResponse)
		assert.NotNil(t, authResponse.Token)
		assert.Equal(t, authResponse.Email, inputtedUser.Email)
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
		user := models.LoginUser{
			Email:    "johndoe@example.com",
			Password: "password",
		}

		// Mock the document that should be returned by FindOne
		hashed, err := HashPassword("password")
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
