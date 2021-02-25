package massdriver

import (
	"encoding/json"
	"reflect"
	"testing"
)



func TestGetConnection(t *testing.T) {
  type test struct {
  	name  string
  	input string
  	want  string
  }

	tests := []test{
		{
			name:  "scalars",
			input: "connId",
			want:  `{"field1": "value1", "field2": "value2"}`
		},
	}

  for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := GetConnection(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestListConnections(t *testing.T) {
  type test struct {
  	name  string
  	input string
  	want  string
  }

	tests := []test{
		{
			name:  "A single variable",
			input: "orgId",
			want:  `["connId1", "connId2": "connId3"]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ListConnections(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}
