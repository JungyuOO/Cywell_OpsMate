package appserver

import (
	"testing"
)

func TestNewUUIDReturnsRFC4122Shape(t *testing.T) {
	id, err := newUUID()
	if err != nil {
		t.Fatal(err)
	}
	if len(id) != 36 {
		t.Fatalf("uuid len = %d, want 36", len(id))
	}
	if id[14] != '4' {
		t.Fatalf("uuid version = %c, want 4", id[14])
	}
}
