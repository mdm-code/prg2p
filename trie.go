package prg2p

import (
	"strings"
)

// trieNode represents double-root trie tree structure holding left/right
// context.
type trieNode struct {
	left, right map[string]*trieNode
	output      []string
	nchars      int
}

// traverseLeft traverses the left context of the trieNode. This method can
// both create the context of the trie or traverse down the existing context.
func (t *trieNode) traverseLeft(s string) *trieNode {
	curr := t
	s = reverse(s)
	for _, c := range strings.Split(s, "") {
		if _, ok := curr.left[c]; !ok {
			t := &trieNode{
				left:  make(map[string]*trieNode),
				right: make(map[string]*trieNode),
			}
			curr.left[c] = t
		}
		curr = curr.left[c]
	}
	return curr
}

// traverseRight traverses the right context of the trieNode. This method can
// both create the context of the trie or traverse down the existing context.
func (t *trieNode) traverseRight(s string) *trieNode {
	curr := t
	for _, c := range strings.Split(s, "") {
		if _, ok := curr.right[c]; !ok {
			t := &trieNode{
				left:  make(map[string]*trieNode),
				right: make(map[string]*trieNode),
			}
			curr.right[c] = t
		}
		curr = curr.right[c]
	}
	return curr
}

// setOutput sets the character count of source and the output word.
func (t *trieNode) setOutput(nchars int, out []string) {
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

// newTree creates a new tree structure with parsed out G2P rules. A tree is
// organized in such a way that top level has characters from the source word,
// tier two holds its right- hand-side context and, tier three its left
// context. The transcription can be figured out based on the top-down
// structure starting at the character and going to the right and then left
// context.
func newTree(i *interpreter) *trieNode {
	t := &trieNode{
		left:  make(map[string]*trieNode),
		right: make(map[string]*trieNode),
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
