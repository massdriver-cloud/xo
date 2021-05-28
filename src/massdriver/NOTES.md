```go
package massdriver

import (
	fmt "fmt"
	"testing"

	mocks "xo/utils/mocks"

	structpb "google.golang.org/protobuf/types/known/structpb"
)

func init() {
	Client = &mocks.MockClient{}
}

func TestGetDeployment(t *testing.T) {
	m, err := structpb.NewValue(map[string]interface{}{
		"firstName": "John",
		"lastName":  "Smith",
		"isAlive":   true,
		"age":       27,
		"address": map[string]interface{}{
			"streetAddress": "21 2nd Street",
			"city":          "New York",
			"state":         "NY",
			"postalCode":    "10021-3100",
		},
		"phoneNumbers": []interface{}{
			map[string]interface{}{
				"type":   "home",
				"number": "212 555-1234",
			},
			map[string]interface{}{
				"type":   "office",
				"number": "646 555-4567",
			},
		},
		"children": []interface{}{},
		"spouse":   nil,
	})

	_ = err

	// log("got %v", deployment.Id)
	// deployment.Connections.AsInterface() > connection.json
	// deployment.Params.AsInterface() > params.json
	json := m.AsInterface()
	fmt.Printf("%v", json)

	artifactsJson := json.Unmarshal(json)
	a, err := structpb.NewValue(artifactsJson)
	massdriver.UploadArtifacts{
		artifacts: a
	}

	got := "lol"
	want := "foo"

	if got != want {
		t.Fatalf("expected: %v, got: %v", want, got)
	}
}
```