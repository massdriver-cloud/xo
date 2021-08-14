package massdriver

import (
	"net/http"
)

type MockHTTPClient struct{}

var (
	MockDoFunc func(req *http.Request) (*http.Response, error)
)

func (m MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return MockDoFunc(req)
}

func init() {
	Client = MockHTTPClient{}
}
