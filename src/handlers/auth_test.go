package handlers

import (
	"gocker-api/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestAuth(t *testing.T) {
	var tests = []struct {
		endpoint     string
		method       string
		expectedCode int
		body         io.Reader
		handler      utils.APIFunc
	}{
		// test registering a correct user
		{"/api/v1/auth/register", "POST", 201,
			strings.NewReader(`{"first_name": "test", "email": "testauth@gmail.com", "password": "testpass"}`),
			handleRegisterUser,
		},
		// test registering an incorrect user
		{"/api/v1/auth/register", "POST", 400,
			strings.NewReader(`{"first_name": "test"}`),
			handleRegisterUser,
		},
		// test authenticating an existing user
		{"/api/v1/auth/authenticate", "POST", 200,
			strings.NewReader(`{"email": "testauth@gmail.com", "password": "testpass"}`),
			handleAuthenticateUser,
		},
		// test authenticating a user with wrong password
		{"/api/v1/auth/authenticate", "POST", 500,
			strings.NewReader(`{"email": "testauth@gmail.com", "password": "wrongpass"}`),
			handleAuthenticateUser,
		},
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

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(utils.ParseToHandlerFunc(test.handler))

		handler.ServeHTTP(rr, req)

		if rr.Code != test.expectedCode {
			t.Errorf("wrong status code. expected %d and got %d, with error %s", test.expectedCode, rr.Code, rr.Body.String())
		}

	}

}
