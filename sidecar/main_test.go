package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRequest(headers http.Header) *http.Request {
	request := httptest.NewRequest("GET", "http://payments.svc.cluster.local:7000/pay", nil)
	request.Header = headers
	return request
}

func Test_isAllowedRequest(t *testing.T) {
	cases := []struct {
		name      string
		request   *http.Request
		isAllowed bool
	}{
		{
			name:      "Valid request",
			request:   newTestRequest(http.Header{"X-Internal-Token": []string{"hidden-token"}}),
			isAllowed: true,
		},
		{
			name:      "Invalid request",
			request:   newTestRequest(http.Header{}),
			isAllowed: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := isAllowedRequest(c.request)
			if result != c.isAllowed {
				t.Errorf("Expected error %v, got %v", c.isAllowed, result)
			}
		})
	}
}
