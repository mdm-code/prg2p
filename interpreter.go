package prg2p

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Variable from an assignment statement with the name (left) and value (right)
// side of the operator.
//
// Examples:
// ALL = a, b, c, d ... n
// SA1 = a, e, y, u, o
type variable struct {
	name  string
	value []string
}

// rule statement with four elements:
// - Left and right context
// - the source letter
// - the target phoneme(s)
//
// Examples:
// PUSTY	ś	-(b)	si
// PUSTY	n	(ni, ci, dzi)	n, ni
// (a, e)	u	PUSTY	l_
type rule struct {
	left   []string
	right  []string
	source string
	target []string
}

// Interpreter interprets G2P rules. It holds two components used to process
// text into phonemic transcription: variables and rules.
type Interpreter struct {
	vars  map[string][]string // Ex. key = ALL, value = a, b, c ... z
	rules []rule
}

// NewInterpreter returns a new Interpreter instance responsible for parsing
// G2P rules. It takes variable assignments and rules as input and creates a
// structure that can be used for building other structures.
func NewInterpreter() *Interpreter {
	i := &Interpreter{
		vars: make(map[string][]string),
	}
	return i
}

// Scan populates Interpreter with G2P rules.
func (i *Interpreter) Scan(f *os.File) error {
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		l = strings.TrimSpace(l)
		if err := i.eval(l); err != nil {
			return fmt.Errorf("could not evaluate %s", l)
		}
	}
	return nil
}

// eval evaluates a line as a variable or a rule.
func (i *Interpreter) eval(l string) error {
	if l == "" || strings.HasPrefix(l, "#") {
		return nil
	}
	if ok := strings.Contains(l, "="); ok {
		err := i.asVar(l)
		if err != nil {
			return err
		}
		return nil
	}
	err := i.asRule(l)
	if err != nil {
		return err
	}
	return nil
}

// asVar evaluates a line as a variable assignment.
func (i *Interpreter) asVar(l string) error {
	sp := strings.Split(l, "=")
	if len(sp) != 2 {
		return fmt.Errorf("multiple assignments on one line %s", l)
	}
	vr := strings.TrimSpace(sp[0])
	vals := strings.Split(sp[1], ",")
	if len(vals) == 1 && strings.TrimSpace(vals[0]) == "" {
		return fmt.Errorf("no values to assign to variable on line %s", l)
	}
	for i := range vals {
		vals[i] = strings.TrimSpace(vals[i])
	}
	if vr == "ALL" {
		vals = append(vals, "$")
	}
	i.vars[vr] = vals
	return nil
}

// asRule evaluates a line as a rule statement.
func (i *Interpreter) asRule(l string) error {
	splits := strings.Split(l, "\t")
	if len(splits) != 4 {
		return fmt.Errorf("expected 4 splits in %s", l)
	}
	lCtx, err := i.context(splits[0])
	if err != nil {
		return err
	}
	rCtx, err := i.context(splits[2])
	if err != nil {
		return err
	}
	var target []string
	for _, s := range strings.Split(splits[3], ",") {
		s = strings.TrimSpace(s)
		target = append(target, s)
	}
	if lCtx != nil && len(lCtx) == 0 {
		return fmt.Errorf("empty left context in line %s", l)
	}
	if rCtx != nil && len(rCtx) == 0 {
		return fmt.Errorf("empty right context in line %s", l)
	}
	r := rule{
		left:   lCtx,
		right:  rCtx,
		source: splits[1],
		target: target,
	}
	i.rules = append(i.rules, r)
	return nil
}

// context returns the left/right context for the source character.
func (i *Interpreter) context(v string) ([]string, error) {
	if s, ok := i.vars[v]; ok && strings.Join(s, "") == "*" {
		return nil, nil
	}
	if _, ok := i.vars["ALL"]; !ok { // "ALL" is the base slice to trim.
		return nil, fmt.Errorf("variable \"ALL\" not set")
	}
	var out []string
	var err error
	if strings.HasPrefix(v, "-") {
		out, err = i.asDifference(v)
		if err != nil {
			return nil, err
		}
	} else {
		out, err = i.asConstraint(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// asDifference returns a limited set of values by removing unwanted values
// them from ALL.
func (i *Interpreter) asDifference(v string) ([]string, error) {
	if strings.HasPrefix(v, "-(") {
		if !strings.HasSuffix(v, ")") {
			return nil, fmt.Errorf("expected \")\" in line %s", v)
		}
		var toRemove []string
		for _, c := range strings.Split(v[2:len(v)-1], ",") {
			c = strings.TrimSpace(c)
			toRemove = append(toRemove, c)
		}
		out := rm(i.vars["ALL"], toRemove)
		return out, nil
	}
	out := i.vars["ALL"]
	for _, vr := range strings.Split(v, "-")[1:] {
		vals, ok := i.vars[vr]
		if !ok {
			return nil, fmt.Errorf("variable \"%s\" not found", vr)
		}
		out = rm(out, vals)
	}
	return out, nil
}

// asConstraint returns a limited set of values.
func (i *Interpreter) asConstraint(v string) ([]string, error) {
	if strings.HasPrefix(v, "(") {
		if !strings.HasSuffix(v, ")") {
			return nil, fmt.Errorf("expected \")\" in line %s", v)
		}
		var toKeep []string
		for _, c := range strings.Split(v[1:len(v)-1], ",") {
			c = strings.TrimSpace(c)
			toKeep = append(toKeep, c)
		}
		return toKeep, nil
	}
	var out []string
	for _, vr := range strings.Split(v, "+") {
		vals, ok := i.vars[vr]
		if !ok {
			return nil, fmt.Errorf("variable \"%s\" not found", vr)
		}
		for _, c := range vals {
			out = append(out, c)
		}
	}
	return out, nil
}

// rm removes items from the first slice if present in the second slice.
func rm(s1, s2 []string) []string {
	var out []string
	ref := make(map[string]bool)
	for _, elem := range s2 {
		ref[elem] = true
	}
	for _, elem := range s1 {
		if _, ok := ref[elem]; !ok {
			out = append(out, elem)
		}
	}
	return out
}
