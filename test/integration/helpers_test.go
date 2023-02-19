package integration_test

import (
	"net/http"
	"testing"

	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/levigross/grequests"
	"github.com/stretchr/testify/require"
)

func doReq(
	t *testing.T,
	method, endpoint string,
	opts *grequests.RequestOptions,
) *grequests.Response {
	t.Helper()

	resp, err := grequests.DoRegularRequest(method, "http://127.0.0.1:8080/"+endpoint, opts)
	require.NoError(t, err)

	return resp
}

func doGetReq(
	t *testing.T,
	endpoint string,
	opts *grequests.RequestOptions,
) *grequests.Response {
	t.Helper()

	return doReq(t, http.MethodGet, endpoint, opts)
}

func doPostReq(
	t *testing.T,
	endpoint string,
	opts *grequests.RequestOptions,
) *grequests.Response {
	t.Helper()

	return doReq(t, http.MethodPost, endpoint, opts)
}

func doLogin(t *testing.T, login, password string) string {
	t.Helper()

	body := schema.LoginRequest{Login: login, Password: password}
	opts := grequests.RequestOptions{JSON: body}

	resp := doPostReq(t, "api/v1/login", &opts)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	token := resp.Header.Get("Authorization")
	require.NotEmpty(t, token)

	return token
}

func doLoginAsHead(t *testing.T) string {
	t.Helper()

	return doLogin(t, "head", "1q2w3e")
}
