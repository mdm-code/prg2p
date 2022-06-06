/*
Package prg2p implements a grapheme-to-phoneme rule-based converter.

Usage
	package main

	import (
		"fmt"

		"github.com/mdm-code/prg2p"
	)

	func main() {
		// Load g2p rules from flat file
		f, err := os.Open("static/rules.txt")
		defer f.Close()
		g2p, err := prg2p.Load(f)

		// Iterate over words to get their phonemic transcripts
		var trans []string
		for _, w := range []string{"ala", "ma", "kota"} {
			t, err := g2p.Transcribe(w, false)
			if err != nil {
				fmt.Println(err)
				continue
			}
			trans = append(trans, t...)
		}
		for _, t := range trans {
			fmt.Println(t)
		}
	}
*/
package prg2p
