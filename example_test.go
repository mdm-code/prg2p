package prg2p_test

import (
	"fmt"

	"github.com/mdm-code/prg2p"
)

func ExampleLoad() {
	r := prg2p.Rules()
	g2p, err := prg2p.Load(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	var trans [][]string
	for _, w := range []string{"ala", "ma", "kta"} {
		t, err := g2p.Transcribe(w, false)
		if err != nil {
			continue
		}
		trans = append(trans, t)
	}
	fmt.Println(trans)
}
