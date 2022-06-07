package prg2p_test

import (
	"fmt"
	"strings"

	"github.com/mdm-code/prg2p"
)

func ExampleNewTree() {
	i := prg2p.NewInterpreter()
	t := prg2p.NewTree(i)
	fmt.Println(t)
}

func ExampleNewInterpreter() {
	i := prg2p.NewInterpreter()
	fmt.Println(i)
}

func ExampleInterpreter() {
	i := prg2p.NewInterpreter()
	r := strings.NewReader(`SA = a, b, c`)
	// Scan rules and declarations from the reader r
	i.Scan(r)
}

func ExampleG2P() {
	g2p := prg2p.NewG2P(nil)
	trans, err := g2p.Transcribe("ala", false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(trans)
}

func ExampleLoad() {
	r := strings.NewReader(`SA = a, b, c`)
	g2p, err := prg2p.Load(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(g2p)
}

func ExampleNewG2P() {
	g2p := prg2p.NewG2P(nil)
	fmt.Println(g2p)
}
