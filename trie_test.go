package prg2p

import "testing"

// Test return empty trieNode on nil pointer.
func TestNilPointer(t *testing.T) {
	tree := NewTree(nil)
	if tree != nil {
		t.Errorf("nil *Interpreter pointer should result in nil tree pointer")
	}
}
