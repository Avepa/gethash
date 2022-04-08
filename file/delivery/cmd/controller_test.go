package cmd

import "testing"

func TestNewController(t *testing.T) {
	ctrl := NewController(nil)
	if ctrl == nil {
		t.Fatal("The NewController should have created an object")
	}
}
