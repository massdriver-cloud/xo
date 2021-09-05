package massdriver

import (
	"io"
	http "net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	mdproto "github.com/massdriver-cloud/rpc-gen-go/massdriver"
)

var (
	Client          mdproto.HTTPClient
	OutputGenerator func(string) (io.Writer, error)
	s               Specification
)

func outputToFile(path string) (io.Writer, error) {
	return os.OpenFile(path, os.O_WRONLY, 0644)
}

func init() {
	Client = &http.Client{}
	OutputGenerator = outputToFile
	envconfig.Process("massdriver", &s)
}

type Specification struct {
	URL string `default:"http://localhost:4000/rpc/workflow"`
}
