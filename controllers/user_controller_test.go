package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
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

		r.ServeHTTP(w, req)
	})
}
