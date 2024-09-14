package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

const Endpoint string = "https://api.massdriver.cloud/api/"

func NewClient(endpoint string, deploymentId string, token string) graphql.Client {
	httpClient := http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport, deploymentId: deploymentId, token: token}}
	return graphql.NewClient(endpoint, &httpClient)
}

type authedTransport struct {
	wrapped      http.RoundTripper
	deploymentId string
	token        string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(t.deploymentId + ":" + t.token))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))
	return t.wrapped.RoundTrip(req)
}
