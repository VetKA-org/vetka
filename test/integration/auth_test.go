package integration_test

import (
	"net/http"
	"testing"

	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/levigross/grequests"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
}

func (ts *AuthTestSuite) TestDoRequestsWithoutToken() {
	endpoints := [...]string{
		"api/v1/appointments",
		"api/v1/patients",
		"api/v1/queue",
		"api/v1/roles",
		"api/v1/species",
		"api/v1/users",
	}

	for _, endpoint := range endpoints {
		ts.T().Run("Access denied to "+endpoint, func(t *testing.T) {
			resp := doGetReq(t, endpoint, nil)

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func (ts *AuthTestSuite) TestSendRequestsWithBadToken() {
	tt := []struct {
		name  string
		token string
	}{
		{
			name:  "Access denied if spaces passed insted of token",
			token: "    ",
		},
		{
			name:  "Access denied if token is malformed",
			token: "xxxxxx!",
		},
		{
			name:  "Access denied if token is signed with different secret",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5NzY0MTE3MzYsImlhdCI6MTY3NjQxMDgzNiwiaXNzIjoiVmV0a2EiLCJqdGkiOiI2ODgxYzY2ZS03ZTNhLTQ1OTAtOWNmZS1lOWJhMDRlODkxZGMiLCJsb2dpbiI6InRlc3QiLCJuYmYiOjE2NzY0MTA4MzYsInJvbGVzIjpbImRvY3RvciJdLCJzdWIiOiJmZThjZWFjOC0xMmRmLTQ0MzktYTM2ZC05YjQ3NDRhZmFmOGUifQ.S-cpHJ0aM8Bp4V943qAOxiQyI4SY0_l-yGpZxhciYYk", //nolint:lll // jwt token
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			headers := map[string]string{"Authorization": "Bearer " + tc.token}
			opts := grequests.RequestOptions{Headers: headers}

			resp := doGetReq(t, "api/v1/users", &opts)

			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	}
}

func (ts *AuthTestSuite) TestBadLogin() {
	tt := []struct {
		name     string
		login    string
		password string
		expected int
	}{
		{
			name:     "Bad request on empty login",
			login:    "",
			password: "1q2w3e",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request on empty password",
			login:    "head",
			password: "",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Unauthorized on unknown login",
			login:    "xxx",
			password: "1q2w3e",
			expected: http.StatusUnauthorized,
		},
		{
			name:     "Unauthorized on unknown login",
			login:    "head",
			password: "123",
			expected: http.StatusUnauthorized,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			req := schema.LoginRequest{Login: tc.login, Password: tc.password}
			opts := grequests.RequestOptions{JSON: req}

			resp := doPostReq(t, "api/v1/login", &opts)

			require.Equal(t, tc.expected, resp.StatusCode)
		})
	}
}
