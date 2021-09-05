package massdriver

import (
	http "net/http"

	"github.com/kelseyhightower/envconfig"
	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
)

var (
	Client mdproto.HTTPClient
	s      Specification
)

func init() {
	Client = &http.Client{}
	envconfig.Process("massdriver", &s)
}

type Specification struct {
	URL string `default:"http://localhost:4000/rpc/workflow"`
}
