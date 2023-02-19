package integration_test

import (
	"net/http"
	"testing"

	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/levigross/grequests"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UsersTestSuite struct {
	suite.Suite

	token string
}

func (ts *UsersTestSuite) SetupTest() {
	ts.token = doLoginAsHead(ts.T())
}

func (ts *UsersTestSuite) TestCreateUserWithBadData() {
	tt := []struct {
		name     string
		login    string
		password string
		roles    []uuid.UUID
		expected int
	}{
		{
			name:     "Bad request if login is empty",
			password: "test-pwd",
			roles:    make([]uuid.UUID, 0),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request if password is empty",
			login:    "test-doctor",
			roles:    make([]uuid.UUID, 0),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request if roles are missing",
			login:    "test-doctor",
			password: "test-pwd",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request if login is too long",
			login:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", //nolint:lll //test data
			password: "test-pwd",
			roles:    make([]uuid.UUID, 0),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request on unknown role id",
			login:    "test-doctor",
			password: "test-password",
			roles:    []uuid.UUID{uuid.NewV4()},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Conflicts if user exists",
			login:    "head",
			password: "test-password",
			roles:    make([]uuid.UUID, 0),
			expected: http.StatusConflict,
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			headers := map[string]string{"Authorization": ts.token}
			body := schema.RegisterUserRequest{
				Login:    tc.login,
				Password: tc.password,
				Roles:    tc.roles,
			}
			opts := grequests.RequestOptions{Headers: headers, JSON: body}

			resp := doPostReq(t, "api/v1/users", &opts)

			require.Equal(t, tc.expected, resp.StatusCode)
		})
	}
}

func (ts *UsersTestSuite) TestCreateUserWithoutRoles() {
	headers := map[string]string{"Authorization": ts.token}
	body := schema.RegisterUserRequest{
		Login:    "guest",
		Password: "guest-password",
		Roles:    make([]uuid.UUID, 0),
	}
	opts := grequests.RequestOptions{Headers: headers, JSON: body}

	resp := doPostReq(ts.T(), "api/v1/users", &opts)

	require.Equal(ts.T(), http.StatusOK, resp.StatusCode)
}
