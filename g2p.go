package prg2p

import (
	"fmt"
	"io"
	"strings"
)

// G2P transcriber class that takes a populated double trie tree
// with parsed grapheme-to-phoneme rules. It exposes transcription
// interface that takes individual words and outputs their most
// likely transcripts.
type G2P struct {
	tree *TrieNode
}

// NewG2P returns G2P object responsibile for handling transcription.
func NewG2P(t *TrieNode) *G2P {
	g := G2P{
		tree: t,
	}
	return &g
}

// Load returns a fully initialized G2P object with rules read from r.
func Load(r io.Reader) (*G2P, error) {
	interp := NewInterpreter()
	err := interp.Scan(r)
	if err != nil {
		return nil, err
	}
	tree := NewTree(interp)
	g2p := NewG2P(tree)
	return g2p, nil
}

// Transcribe word from graphemic to phonemic transcription. Use n to specify
// whether to return all possible transcriptions or just the first hit.
func (g *G2P) Transcribe(w string, all bool) ([]string, error) {
	if g.tree == nil {
		return []string{}, fmt.Errorf("tire node is nil")
	}
	var trans [][]string
	w = strings.ToLower(w)
	nchars := len([]rune(w))
	i := 0
	for i < nchars {
		t := g.rightVars(w, i, i-1, g.tree)
		if t == nil {
			return []string{}, fmt.Errorf("failed to transcribe %s", w)
		}
		trans = append(trans, t.output)
		i += t.nchars
	}
	out, err := g.all(trans, 0)
	if err != nil {
		return []string{}, err
	}
	if all == true {
		return out, nil
	}
	return out[:1], nil
}

// All grabs all possible transcription variants.
func (g *G2P) all(trans [][]string, i int) ([]string, error) {
	if len(trans) == 0 {
		return []string{}, fmt.Errorf("no transcription variants offered")
	}
	if i == len(trans)-1 {
		return trans[len(trans)-1], nil
	}
	rest, err := g.all(trans, i+1)
	if err != nil {
		return []string{}, err
	}
	var result []string
	for _, i := range trans[i] {
		for _, j := range rest {
			result = append(result, i+" "+j)
		}
	}
	return result, nil
}

// RightVars traverses the right-hand side of the complete double trie.
func (g *G2P) rightVars(w string, frontIdx, backIdx int, trie *TrieNode) *TrieNode {
	wRune := []rune(w)
	var curChar string
	if frontIdx < len(wRune) {
		curChar = string(wRune[frontIdx])
	}
	if t, ok := trie.right[curChar]; frontIdx < len(wRune) && ok {
		frontIdx++
		t := g.rightVars(w, frontIdx, backIdx, t)
		if t != nil {
			return t
		}
	}
	if t, ok := trie.right["$"]; frontIdx == len(wRune) && ok {
		t := g.leftVars(w, backIdx, t)
		if t != nil {
			return t
		}
	}
	t := g.leftVars(w, backIdx, trie)
	if t != nil {
		return t
	}
	if trie.nchars != 0 {
		return trie
	}
	return nil
}

// LeftVars traverses left-hand side part of the complete double trie.
func (g *G2P) leftVars(w string, backIdx int, trie *TrieNode) *TrieNode {
	wRune := []rune(w)
	curChar := string(wRune[len(wRune)-2-backIdx])

	if t, ok := trie.left[curChar]; backIdx >= 0 && ok {
		backIdx--
		t := g.leftVars(w, backIdx, t)
		if t != nil {
			return t
		}
	}
	if t, ok := trie.left["$"]; backIdx == -1 && ok {
		if t.nchars != 0 {
			return t
		}
	}
	if trie.nchars != 0 {
		return trie
	}
	return nil
}
