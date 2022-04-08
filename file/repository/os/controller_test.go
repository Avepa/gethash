package os

import "testing"

func TestNewController(t *testing.T) {
	ctrl := NewController()
	if ctrl == nil {
		t.Fatal("The NewController should have created an object")
	}
}
