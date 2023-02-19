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

func getRole(t *testing.T, token, name string) uuid.UUID {
	t.Helper()

	headers := map[string]string{"Authorization": token}
	opts := grequests.RequestOptions{Headers: headers}

	resp := doGetReq(t, "api/v1/roles?name="+name, &opts)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := schema.ListRolesResponse{}
	err := resp.JSON(&body)
	require.NoError(t, err)
	require.Equal(t, 1, len(body.Data))

	return body.Data[0].ID
}

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
			login:    "test",
			roles:    make([]uuid.UUID, 0),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request if roles are missing",
			login:    "test",
			password: "pwd",
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request if login is too long",
			login:    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", //nolint:lll //test data
			password: "pwd",
			roles:    make([]uuid.UUID, 0),
			expected: http.StatusBadRequest,
		},
		{
			name:     "Bad request on unknown role id",
			login:    "test",
			password: "pwd",
			roles:    []uuid.UUID{uuid.NewV4()},
			expected: http.StatusBadRequest,
		},
		{
			name:     "Conflicts if user exists",
			login:    "head",
			password: "pwd",
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

func (ts *UsersTestSuite) TestCreateUser() {
	tt := []struct {
		name     string
		login    string
		password string
		roleName string
	}{
		{
			name:     "Create user with head doctor role",
			login:    "another-head",
			password: "head-pwd",
			roleName: "head",
		},
		{
			name:     "Create user with doctor role",
			login:    "test-doctor",
			password: "doctor-pwd",
			roleName: "doctor",
		},
		{
			name:     "Create user with administrator role",
			login:    "test-admin",
			password: "admin-pwd",
			roleName: "admin",
		},
	}

	for _, tc := range tt {
		ts.T().Run(tc.name, func(t *testing.T) {
			headers := map[string]string{"Authorization": ts.token}

			roleID := getRole(t, ts.token, tc.roleName)

			body := schema.RegisterUserRequest{
				Login:    tc.login,
				Password: tc.password,
				Roles:    []uuid.UUID{roleID},
			}
			opts := grequests.RequestOptions{Headers: headers, JSON: body}

			resp := doPostReq(t, "api/v1/users", &opts)

			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func (ts *UsersTestSuite) TestCreateUserWithSeveralRoels() {
	t := ts.T()

	headers := map[string]string{"Authorization": ts.token}

	headRoleID := getRole(t, ts.token, "head")
	doctorRoleID := getRole(t, ts.token, "doctor")

	body := schema.RegisterUserRequest{
		Login:    "multi-roles",
		Password: "pwd",
		Roles:    []uuid.UUID{headRoleID, doctorRoleID},
	}
	opts := grequests.RequestOptions{Headers: headers, JSON: body}

	resp := doPostReq(t, "api/v1/users", &opts)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func (ts *UsersTestSuite) TestListUsers() {
	t := ts.T()

	headers := map[string]string{"Authorization": ts.token}
	opts := grequests.RequestOptions{Headers: headers}

	resp := doGetReq(t, "api/v1/users", &opts)

	require.Equal(ts.T(), http.StatusOK, resp.StatusCode)

	body := schema.ListUsersResponse{}
	err := resp.JSON(&body)

	require.NoError(t, err)
	require.NotZero(t, len(body.Data))

	for _, user := range body.Data {
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.Login)
		require.Empty(t, user.Password)
	}
}
