package massdriver

import (
	http "net/http"

	"github.com/kelseyhightower/envconfig"
)

var (
	Client HTTPClient
	s      Specification
)

func init() {
	Client = &http.Client{}
	envconfig.Process("massdriver", &s)
}

type Specification struct {
	URL string `default:"http://localhost:4000/rpc/workflow"`
}
