package prg2p

import (
	"testing"
)

// Test if Rules returns an object impmementing the io.Reader interface.
func TestRulesCall(t *testing.T) {
	buf := []byte{}
	r := Rules()
	r.Read(buf)
}
