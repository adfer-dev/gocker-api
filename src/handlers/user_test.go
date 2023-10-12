package handlers

import (
	"gocker-api/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestUsers(t *testing.T) {
	var tests = []struct {
		endpoint     string
		method       string
		expectedCode int
		param        string
		body         io.Reader
		handler      utils.APIFunc
	}{
		// test getting all users
		{"/api/v1/users", "GET", 200, "", nil, handleGetUsers},
		// test getting an specific user
		{"/api/v1/users/{id}", "GET", 200, "1", nil, handleGetUser},
		// test getting a not existent user
		{"/api/v1/users/{id}", "GET", 404, "10000", nil, handleGetUser},
		// test adding a correct user
		{"/api/v1/users", "POST", 201, "",
			strings.NewReader(`{"first_name": "test", "email": "test@gmail.com", "password": "testpass"}`),
			handleCreateUser,
		},
		// test adding an incorrect user
		{"/api/v1/users", "POST", 400, "",
			strings.NewReader(`{"first_name": "test"}`),
			handleCreateUser,
		},
		// test updating an existing user
		{"/api/v1/users/{id}", "PUT", 201, "10",
			strings.NewReader(`{"first_name": "updatedtest"}`),
			handleUpdateUser,
		},
		// test deleting a user
		{"/api/v1/users/{id}", "DELETE", 201, "10", nil, handleDeleteUser},
	}

	err := godotenv.Load("../.env")

	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		req, reqErr := http.NewRequest(test.method, test.endpoint, test.body)

		if reqErr != nil {
			t.Fatal(reqErr)
		}

		req.Header.Add("Authorization", os.Getenv("TEST_TOKEN"))

		if test.param != "" {
			req = mux.SetURLVars(req, map[string]string{"id": test.param})
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(utils.ParseToHandlerFunc(test.handler))

		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedCode {
			t.Errorf("wrong status code. expected %d and got %d, with error %s", test.expectedCode, rr.Code, rr.Body.String())
		}

	}
}
