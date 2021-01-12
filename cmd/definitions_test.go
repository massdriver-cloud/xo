package cmd

import (
	"testing"
)

func TestFoo(t *testing.T) {
	got := true
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}
