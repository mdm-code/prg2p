package prg2p

import (
	"io"
	"reflect"
	"testing"
)

// Fresh instance of populated *TrieNode for unit testing.
func inputTrie() (*TrieNode, error) {
	i := NewInterpreter()
	err := i.Scan(rulesIO())
	if err != nil {
		return nil, err
	}
	t := NewTree(i)
	return t, nil
}

// Test the Load package interface function of prg2p.
func TestLoad(t *testing.T) {
	_, err := Load(rulesIO())
	if err != nil {
		t.Errorf("failed to load G2P rules transcriber")
	}
}

// Test if passing nil io.Reader interface causes Load to error out.
func TestLoadFails(t *testing.T) {
	var r io.Reader
	_, err := Load(r)
	if err == nil {
		t.Errorf("nil interface should cause an error")
	}
}

// Test empty Transcribe errors out when trie is nil
func TestErrorTrieNil(t *testing.T) {
	g2p := NewG2P(nil)
	_, err := g2p.Transcribe("test", false)
	if err == nil {
		t.Error("nil tree should cause Transcribe to fail")
	}
}

// Test G2P.Transcribe() package interface method.
func TestTranscribe(t *testing.T) {
	cases := []struct {
		name string
		word string
		all  bool
		want []string
	}{
		{"ala-f", "ala", false, []string{"a l a"}},
		{"ma-f", "ma", false, []string{"m a"}},
		{"ma-t", "ma", true, []string{"m a"}},
		{"kota-f", "Kota", false, []string{"k o t a"}},
		{"kota-f-capital", "Kota", false, []string{"k o t a"}},
		{"chcę-t-capital", "Chcę", true, []string{"h c e", "h c e_"}},
		{"mówię-t", "mówię", true, []string{"m u w i e", "m u w i e_"}},
	}
	g2p, err := Load(rulesIO())
	if err != nil {
		t.Fatal("failed to create G2P transcriber")
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			have, err := g2p.Transcribe(c.word, c.all)
			if err != nil {
				t.Errorf("failed to transcribe: %s", c.word)
			}
			if ok := reflect.DeepEqual(have, c.want); !ok {
				t.Errorf("have %v; want: %v", have, c.want)
			}
		})
	}
}

// Test if the G2P.Transcribe() package interface method fails gracefully.
func TestTranscribeFails(t *testing.T) {
	cases := []struct {
		name, word string
		all        bool
	}{
		{"number", "5432", false},
		{"punctuation", "wiedzie,", false},
		{"space", "i tak", false},
		{"space-punctuation", "i tak...", false},
		{"empty", "", false},
	}
	g2p, err := Load(rulesIO())
	if err != nil {
		t.Fatal("failed to create G2P transcriber")
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := g2p.Transcribe(c.word, c.all)
			if err == nil {
				t.Errorf("word %s was expected to cause error", c.word)
			}
		})
	}
}
