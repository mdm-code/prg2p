package prg2p

import (
	"strings"
)

// TrieNode represents double-root trie tree structure holding left/right
// context.
type TrieNode struct {
	left, right map[string]*TrieNode
	output      []string
	nchars      int
}

// traverseLeft traverses the left context of the trieNode. This method can
// both create the context of the trie or traverse down the existing context.
func (t *TrieNode) traverseLeft(s string) *TrieNode {
	curr := t
	s = reverse(s)
	for _, c := range strings.Split(s, "") {
		if _, ok := curr.left[c]; !ok {
			t := &TrieNode{
				left:  make(map[string]*TrieNode),
				right: make(map[string]*TrieNode),
			}
			curr.left[c] = t
		}
		curr = curr.left[c]
	}
	return curr
}

// traverseRight traverses the right context of the trieNode. This method can
// both create the context of the trie or traverse down the existing context.
func (t *TrieNode) traverseRight(s string) *TrieNode {
	curr := t
	for _, c := range strings.Split(s, "") {
		if _, ok := curr.right[c]; !ok {
			t := &TrieNode{
				left:  make(map[string]*TrieNode),
				right: make(map[string]*TrieNode),
			}
			curr.right[c] = t
		}
		curr = curr.right[c]
	}
	return curr
}

// setOutput sets the character count of source and the output word.
func (t *TrieNode) setOutput(nchars int, out []string) {
	if t.nchars > nchars {
		return
	}
	var i int
	i, t.nchars, t.output = t.nchars, nchars, out
	if i > 0 {
		return
	}
}

// reverse a string.
func reverse(s string) string {
	var out string
	for _, c := range s {
		out = string(c) + out
	}
	return out
}

// NewTree creates a new tree structure with parsed out G2P rules. A tree is
// organized in such a way that top level has characters from the source word,
// tier two holds its right- hand-side context and, tier three its left
// context. The transcription can be figured out based on the top-down
// structure starting at the character and going to the right and then left
// context.
func NewTree(i *Interpreter) *TrieNode {
	t := &TrieNode{
		left:  make(map[string]*TrieNode),
		right: make(map[string]*TrieNode),
	}
	if i == nil {
		return nil
	}
	for _, r := range i.rules {
		l, r, src, tgt := r.left, r.right, r.source, r.target
		tierOne := t.traverseRight(src)
		if l == nil {
			l = []string{""}
		}
		if r == nil {
			r = []string{""}
		}
		for _, tkn := range r {
			tierTwo := tierOne.traverseRight(tkn)
			for _, tkn := range l {
				tierThree := tierTwo.traverseLeft(tkn)
				tierThree.setOutput(len([]rune(src)), tgt)
			}
		}
	}
	return t
}
