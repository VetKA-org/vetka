package integration_test

import (
	"net/http"
	"testing"

	"github.com/levigross/grequests"
	"github.com/stretchr/testify/require"
)

func sendReq(
	t *testing.T,
	method, endpoint string,
	opts *grequests.RequestOptions,
) *grequests.Response {
	t.Helper()

	resp, err := grequests.DoRegularRequest(method, "http://127.0.0.1:8080/"+endpoint, opts)
	require.NoError(t, err)

	return resp
}

func sendGetReq(
	t *testing.T,
	endpoint string,
	opts *grequests.RequestOptions,
) *grequests.Response {
	t.Helper()

	return sendReq(t, http.MethodGet, endpoint, opts)
}

func sendPostReq(
	t *testing.T,
	endpoint string,
	opts *grequests.RequestOptions,
) *grequests.Response {
	t.Helper()

	return sendReq(t, http.MethodPost, endpoint, opts)
}
