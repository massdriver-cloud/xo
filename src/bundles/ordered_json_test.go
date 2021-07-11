package bundles

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestMarshalJSON(t *testing.T) {
	type test struct {
		name  string
		input string
		want  string
	}

	tests := []test{
		{
			name: "normal test",
			input: `foo: bar
name: John Doe
address:
    street: 123 E 3rd St
    city: Denver
    test: ["a", "b", "c"]
    state: CO
    zip: 81526
anotherTest: [1, 2, 3, 4]`,
			want: `{"foo":"bar","name":"John Doe","address":{"street":"123 E 3rd St","city":"Denver","test":["a","b","c"],"state":"CO","zip":81526},"anotherTest":[1,2,3,4]}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oj := OrderedJSON{}
			err := yaml.Unmarshal([]byte(tc.input), &oj)
			if err != nil {
				println(err)
			}

			bytes, _ := json.Marshal(oj)
			got := string(bytes)

			if got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}
