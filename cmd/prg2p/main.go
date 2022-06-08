package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mdm-code/prg2p"
)

var (
	rule string
	all  bool
)

const (
	exitSuccess = iota
	exitFailure
)

const usage = `prg2p - grapheme-to-phoneme converter

The prg2p utility reads space-delimited words sequentially from standard input,
writing converted phonemic transcripts to standard output.

Usage:   prg2p [-h] [-r FILE] [-a BOOL] [FILE ...]

Example: echo ala ma kota | prg2p -r=rules.txt -a=false

Options:
	-h, --help  show this help message and exit
	-r, --rule  file with g2p rules (default: prg2p.Rules())
	-a, --all   print all allowed conversions (default: False)
`

func main() {
	flag.StringVar(&rule, "r", "", "")
	flag.StringVar(&rule, "rules", "", "")
	flag.BoolVar(&all, "a", false, "")
	flag.BoolVar(&all, "all", false, "")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	var (
		f   io.Reader
		err error
	)
	if rule == "" {
		f = prg2p.Rules()
	} else {
		f, err = os.Open(rule)
		if err != nil {
			fmt.Fprintf(os.Stderr, EOL(err.Error()))
			os.Exit(exitFailure)
		}
	}

	g2p, err := prg2p.Load(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, EOL(err.Error()))
		os.Exit(exitFailure)
	}

	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)
	out := bufio.NewWriter(os.Stdout)

	for in.Scan() {
		word := in.Text()
		trans, err := g2p.Transcribe(word, all)
		if err != nil {
			fmt.Fprintf(os.Stderr, EOL(err.Error()))
			os.Exit(exitFailure)
		}
		line := FTrans(word, trans)
		_, err = out.Write([]byte(EOL(line)))
		if err != nil {
			fmt.Fprintf(os.Stderr, EOL(err.Error()))
			os.Exit(exitFailure)
		}
	}
	if err := out.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, EOL(err.Error()))
		os.Exit(exitFailure)
	}
	os.Exit(exitSuccess)
}

// EOL returns the string s with newline character at the end.
func EOL(s string) string {
	return s + "\n"
}

// FTrans collates a single converted output line.
func FTrans(word string, trans []string) string {
	joined := strings.Join(trans, "|")
	return word + "\t" + strconv.Itoa(len(trans)) + "\t" + joined
}
