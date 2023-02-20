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

func newRegisterUserRequest() schema.RegisterUserRequest {
	return schema.RegisterUserRequest{
		Login:    "test",
		Password: "test-pwd",
		Roles:    make([]uuid.UUID, 0),
	}
}

type UsersTestSuite struct {
	suite.Suite

	token string
}

func (ts *UsersTestSuite) SetupTest() {
	ts.token = doLoginAsHead(ts.T())
}

func (ts *UsersTestSuite) TestRegisterUserWithoutRequiredFields() {
	require := ts.Require()

	headers := map[string]string{"Authorization": ts.token}
	req := schema.RegisterUserRequest{}
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(ts.T(), "api/v1/users", &opts)
	require.Equal(http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(ts.T(), resp)
	expected := []string{"login", "password", "roles"}
	require.Equal(expected, failures.Fields())
}

func (ts *UsersTestSuite) TestRegisterUserWithBadLogin() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterUserRequest()
	req.Login = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" //nolint:lll //test data

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/users", &opts)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	failures := extractErrors(t, resp)
	require.Len(t, failures.Errors, 1)
	require.Equal(t, "login", failures.Errors[0].Field)
}

func (ts *UsersTestSuite) TestRegisterUserWithUnknownRole() {
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterUserRequest()
	req.Roles = []uuid.UUID{uuid.NewV4()}

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(ts.T(), "api/v1/users", &opts)
	require.Equal(ts.T(), http.StatusBadRequest, resp.StatusCode)
}

func (ts *UsersTestSuite) TestRegisterUserIfUserExists() {
	headers := map[string]string{"Authorization": ts.token}

	req := newRegisterUserRequest()
	req.Login = "head"

	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(ts.T(), "api/v1/users", &opts)
	require.Equal(ts.T(), http.StatusConflict, resp.StatusCode)
}

func (ts *UsersTestSuite) TestCreateUserWithoutRoles() {
	t := ts.T()
	headers := map[string]string{"Authorization": ts.token}

	req := schema.RegisterUserRequest{
		Login:    "guest",
		Password: "guest-password",
		Roles:    make([]uuid.UUID, 0),
	}
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/users", &opts)
	require.Equal(t, http.StatusOK, resp.StatusCode)
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

			req := schema.RegisterUserRequest{
				Login:    tc.login,
				Password: tc.password,
				Roles:    []uuid.UUID{roleID},
			}
			opts := grequests.RequestOptions{Headers: headers, JSON: req}

			resp := doPostReq(t, "api/v1/users", &opts)

			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func (ts *UsersTestSuite) TestCreateUserWithSeveralRoles() {
	t := ts.T()

	headers := map[string]string{"Authorization": ts.token}

	headRoleID := getRole(t, ts.token, "head")
	doctorRoleID := getRole(t, ts.token, "doctor")

	req := schema.RegisterUserRequest{
		Login:    "multi-roles",
		Password: "pwd",
		Roles:    []uuid.UUID{headRoleID, doctorRoleID},
	}
	opts := grequests.RequestOptions{Headers: headers, JSON: req}

	resp := doPostReq(t, "api/v1/users", &opts)

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func (ts *UsersTestSuite) TestListUsers() {
	require := ts.Require()

	headers := map[string]string{"Authorization": ts.token}
	opts := grequests.RequestOptions{Headers: headers}

	resp := doGetReq(ts.T(), "api/v1/users", &opts)

	require.Equal(http.StatusOK, resp.StatusCode)

	body := schema.ListUsersResponse{}
	err := resp.JSON(&body)

	require.NoError(err)
	require.NotZero(len(body.Data))

	for _, user := range body.Data {
		require.NotEmpty(user.ID)
		require.NotEmpty(user.Login)
		require.Empty(user.Password)
	}
}
