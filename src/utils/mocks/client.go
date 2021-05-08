package mocks

import (
  "net/http"
)

type MockClient struct {}

var (
	MockDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return MockDoFunc(req)
}
