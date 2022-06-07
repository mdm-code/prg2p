package prg2p

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func rulesIO() *strings.Reader {
	result := strings.NewReader(`
# LITERY
# a, ą, b, c, ć, d, e, ę, f, g, h, i, j, k, l, ł, m, n, ń, o, ó, p, q, r, s, ś, t, u, v, w, x, y, z, ź, ż, é, ü, ö, š, ë

# DWUZNAKI I ZMIĘKCZENIA
# ch, cz, dz, dź, dż, rz, sz, ci, ni, si, zi, dzi

ALL = a, ą, b, c, ć, d, e, ę, f, g, h, i, j, k, l, ł, m, n, ń, o, ó, p, q, r, s, ś, t, u, v, w, x, y, z, ź, ż, ch, cz, dz, dź, dż, rz, sz, ci, ni, si, zi, dzi, é, ü, ö, š, ë, $

# PUSTY - dowolne znaki w kontekście, w tym koniec i początek wyrazu

PUSTY = *

# SA – samogłoski

SA = a, e, y, ą, ę, u, o, i, ó

# SA1 - samogłoski:

SA1 = a, e, y, ą, ę, u, o, ó

# SP – spółgłoski:

SP = b, c, ć, ch, cz, d, dz, dż, f, g, h, j, k, l, ł, m, n, ń, p, r, rz, s, ś, sz, t, w, z, ż, q, v, x, š

# SB - spółgłoski bezdźwięczne:

SB = p, k, t, s, c, h, f, ć, ś, ch, cz, sz, ci

# SD - spółgłoski dźwięczne:

SD = b, d, dz, dź, dż, g, rz, w, z, ź, ż

# SP1 - spółgłoski:

SP1 = p, b, m, w, f, ch, h

# SP2 – spółgłoski:

SP2 = l, r, t, d, k, g

# FR - spółgłoski szczelinowe:

FR = f, w, s, z, ź, ś, ż, h, sz

# SPZ - spółgłoski do zapisu zi

SPZ = d, g, l, n, r, sz, ś

# END - koniec wyrazu

END = $


# GŁOSKI
# a, a_, b, c, ci, cz, d, dz, dzi, drz, e, e_, f, g, h, i, j, k, l, l_, m, n, ni, o, p, r, s, si, sz, t, u, w, y, z, zi, rz
# + ng  + ni_ + h_
# (bank, pański, niechby/Bohdan)

# WYJĄTKI + ODMIANY – do wpisania ręcznie
#		tysiąchektarowy, mierzić, obmierzły, zinterpretować, ziścić, ziszczać, nadziemny, nadziemski, podziemie, podziemny,

# przed	znak	po	głoska

# zawsze

PUSTY	a	PUSTY	a
PUSTY	c	PUSTY	c
PUSTY	e	PUSTY	e
PUSTY	j	PUSTY	j
PUSTY	l	PUSTY	l
PUSTY	ł	PUSTY	l_
PUSTY	m	PUSTY	m
PUSTY	o	PUSTY	o
PUSTY	ó	PUSTY	u
PUSTY	p	PUSTY	p
PUSTY	r	PUSTY	r
END	y	PUSTY	j
-END	y	PUSTY	y

# ubezdźwięcznienie

PUSTY	b	END	p, b
PUSTY	b	SB	p
PUSTY	b	-SB-END	b
PUSTY	dz	END	c, dz
PUSTY	dz	SB	c
PUSTY	dz	-SB-END	dz
PUSTY	dź	END	ci, dzi
PUSTY	dź	SB	ci
PUSTY	dź	-SB-END	dzi
PUSTY	d	END	t, d
PUSTY	d	SB	t
PUSTY	d	-SB-END	d
PUSTY	g	END	k, g
PUSTY	g	SB	k
PUSTY	g	-SB-END	g
PUSTY	rz	END	sz, rz
PUSTY	rz	SB	sz
SB	rz	PUSTY	sz
-SB	rz	-SB-END	rz
PUSTY	w	END	f, w
PUSTY	w	SB	f
SB	w	PUSTY	f
-SB	w	-SB-END	w
PUSTY	z	END	s,z
PUSTY	z	SB	s
PUSTY	z	-SB-END	z
PUSTY	ź	END	si, zi
PUSTY	ź	-END	zi
PUSTY	ż	END	sz, rz
PUSTY	ż	SB	sz
SB	ż	PUSTY	sz
-SB	ż	-SB-END	rz
PUSTY	dż	END	cz, drz
PUSTY	dż	-END	drz

# udźwięcznienie
PUSTY	ć	(b)	ci, dzi
PUSTY	ć	-(b)	ci
PUSTY	ś	(b)	si, zi
PUSTY	ś	-(b)	si
PUSTY	s	(b, g)	s, z
PUSTY	s	-(b, g)	s
PUSTY	t	(b, g)	t, d
PUSTY	t	-(b, g)	t
PUSTY	k	(b, ż)	k, g
PUSTY	k	-(b, ż)	k
PUSTY	cz	(b, d)	cz, drz
PUSTY	cz	-(b, d)	cz
PUSTY	f	(g)	f, w
PUSTY	f	-(g)	f

# inne

PUSTY	h	SD	h
PUSTY	h	-SD	h
SP1	i	SA	j
SP2	i	SA	i, j
SP1	i	-SA	i
SP2	i	-SA	i
-SP1-SP2	i	PUSTY	i
PUSTY	n	(k, g)	n
PUSTY	n	(ni, ci, dzi)	n, ni
PUSTY	n	-(k, g, ni, ci, dzi)	n
PUSTY	ń	FR	ni
PUSTY	ń	-FR	ni
(a, e)	u	PUSTY	l_
-(a, e)	u	PUSTY	u

# dwuznak

PUSTY	sz	PUSTY	sz
PUSTY	ch	SD	h
PUSTY	ch	-SD	h
PUSTY	ci	SA1	ci
PUSTY	ci	SP	ci i
PUSTY	ci	END	ci i
PUSTY	dzi	SA1	dzi
PUSTY	dzi	SP	dzi i
PUSTY	dzi	END	dzi i
PUSTY	ni	SA1	ni
PUSTY	ni	(i)	ni j, ni
PUSTY	ni	SP	ni i
PUSTY	ni	END	ni i
PUSTY	si	SA1	si
PUSTY	si	SP	si i
PUSTY	si	END	si i
PUSTY	zi	SA1	zi
PUSTY	zi	-SA1-SPZ	zi i
PUSTY	zi	SPZ	z i
PUSTY	zi	END	zi i

# ąę

PUSTY	ę	END	e, e_
PUSTY	ą	END	o l_, a_, o m
PUSTY	ę	(l, ł)	e
PUSTY	ą	(ł)	o
PUSTY	ę	(b, p)	e m
PUSTY	ą	(b, p)	o m
PUSTY	ę	(m, n)	e
PUSTY	ą	(m)	o
PUSTY	ą	(n, j)	a_, o l_
PUSTY	ę	(w, f, s, z, ż, rz, sz, ch, h)	e_, e l_
PUSTY	ą	(w, f, s, z, ż, rz, sz, ch, h)	a_, o l_
PUSTY	ę	(d, t, c, dz, dż, cz)	e_, e n
PUSTY	ą	(d, t, c, dz, dż, cz)	a_, o n
PUSTY	ę	(dź, dzi, ć, ci)	e_, e ni
PUSTY	ą	(dź, dzi, ć, ci)	a_, o ni
PUSTY	ę	(g, k)	e n, e_
PUSTY	ą	(g, k)	o n, a_
PUSTY	ę	(ź, zi, ś, si)	e l_, e ni
PUSTY	ą	(ź, zi, ś, si)	o l_, o ni

# warianty, wyjątki, zbitki

PUSTY	trz	PUSTY	t sz, cz
PUSTY	drz	PUSTY	d rz, drz
PUSTY	zsz	PUSTY	s sz, sz
PUSTY	nadz	-(i)	n a d z
PUSTY	nadż	PUSTY	n a d rz
PUSTY	podz	-(i)	p o d z
PUSTY	podż	PUSTY	p o d rz
END	odz	-(i)	o d z
PUSTY	odż	PUSTY	o d rz
PUSTY	budż	PUSTY	b u d rz
PUSTY	śćdzi	PUSTY	si dzi, zi dzi, si ci dzi
PUSTY	śćs	PUSTY	si s, j s, si ci s
PUSTY	śródzi	PUSTY	si r u d zi
PUSTY	sji	PUSTY	s j i, s i
PUSTY	cji	PUSTY	c j i, c i
PUSTY	izm	END	i z m, i s m
PUSTY	izm	(i)	i z m, i zi m
PUSTY	on	(s)	o n, a_
-END	en	(t, k, ci, s)	e n, e_
-(n)	ii	END	i i, i, j i
PUSTY	dźm	(y)	ci m
PUSTY	marz	(ł, n, l)	m a r z
PUSTY	żć	END	si ci, zi ci
PUSTY	klie	PUSTY	k l i j e

# do naprawienia

PUSTY	czw	PUSTY	cz f
PUSTY	szw	PUSTY	sz f

#obce

PUSTY	v	PUSTY	w
PUSTY	q	PUSTY	k
PUSTY	x	PUSTY	k s
PUSTY	é	PUSTY	e
PUSTY	ü	PUSTY	u, i
PUSTY	ö	PUSTY	y
PUSTY	š	PUSTY	s
PUSTY	ë	PUSTY	e`)
	return result
}

// Check if the interpreter handles the file with rules without raising an
// error.
func TestParser(t *testing.T) {
	i := NewInterpreter()
	err := i.Scan(rulesIO())
	if err != nil {
		t.Errorf("%v", err)
	}
}

// Test errScan raised with nil io.Reader interface.
func TestScanning(t *testing.T) {
	i := NewInterpreter()
	var r io.Reader
	err := i.Scan(r)
	if err != nil {
		if errors.Is(err, errScan) {
			return
		}
		t.Errorf("%v", err)
	}
}

// Test Interpreter fails on unexpected line in the reader.
func TestInterpreterFails(t *testing.T) {
	i := NewInterpreter()
	r := strings.NewReader(`she sells sea shells at the sea shore`)
	err := i.Scan(r)
	if err == nil {
		t.Errorf("interpreter Scan() call should fail")
	}
}

// Verify each of the possible outputs for line evaluation.
func TestEval(t *testing.T) {
	i := &Interpreter{
		vars: make(map[string][]string),
	}
	i.vars["PUSTY"] = []string{"*"}
	i.vars["SA1"] = []string{"a", "e", "y", "ą", "ę", "u", "o", "ó"}
	valid := []string{
		"ALL = a, b, c, d, e",
		"SB  = dz, rz, cz, d",
		"PUSTY = *",
		"SA1 = a, g, t, w",
		"# This is a comment.",
		"PUSTY	dz	SA1	dzi",
		"",
	}
	invalid := []string{
		"END	cz	SA2	czi	dzi", // Five elements instead of four
		"ALL = dz, cz = SB", // Multiple assignments on one line
	}
	for _, l := range valid {
		err := i.eval(l)
		if err != nil {
			t.Errorf("error was raised: %s", err)
		}
	}
	for _, l := range invalid {
		err := i.eval(l)
		if err == nil {
			t.Errorf("error was not raised")
		}
	}
}

// Test if rules are parsed into their respective structures based on their
// input.
func TestAsRule(t *testing.T) {
	i := &Interpreter{
		vars: make(map[string][]string),
	}
	i.vars["ALL"] = []string{"a", "e", "i", "o", "u", "y"}
	i.vars["SB"] = []string{"dz", "cz", "sz"}
	i.vars["PUSTY"] = []string{"*"}
	i.vars["SA"] = []string{"r", "n", "l", "m"}
	i.vars["END"] = []string{"$"}

	valid := []string{
		"PUSTY	cz	END	sz",
		"SA	a	SB	e",
		"-SA	a	-SB	y",
		"-(a, u, y)	e	END	i",
		"(i, j)	i	-SA	j",
	}
	invalid := []string{
		"END	cz	END",
		"-(a, u, i)	r	(a, b)	l	END",
		"(a, e	u	END	y",
		"PUSTY	y	-(a, e	u",
	}

	for _, l := range valid {
		err := i.asRule(l)
		if err != nil {
			t.Errorf("error was raised: %s", err)
		}
	}
	for _, l := range invalid {
		err := i.asRule(l)
		if err == nil {
			t.Errorf("error was not raised")
		}
	}
}

// Check if variables are set given the correct input.
func TestAsVariable(t *testing.T) {
	i := &Interpreter{
		vars: make(map[string][]string),
	}
	valid := []string{
		"ALL = a, b, c, d, e",
		"SB  = dz, rz, cz, d",
		"SA=p,t,k",
		"PUSTY = *",
	}
	invalid := []string{
		"ALL = dz, cz = SB",
		"",
		"SB = ",
		"ALL k, g, n",
		"ALL k g n",
	}

	for _, l := range valid {
		err := i.asVar(l)
		if err != nil {
			t.Errorf("error was raised: %s", err)
		}
	}
	for _, l := range invalid {
		err := i.asVar(l)
		if err == nil {
			t.Errorf("error was not raised")
		}
	}
}

// Remove items from slice one if they're present in slice two.
func TestRmFromSlice(t *testing.T) {
	first := []string{"a", "b", "c", "d"}
	second := []string{"c", "d"}
	want := []string{"a", "b"}
	has := rm(first, second)
	if ok := reflect.DeepEqual(has, want); !ok {
		t.Errorf("rmFromSlice want %v; has %v", want, has)
	}
}

// Check if constrained context is created correctly.
func TestAsConstraint(t *testing.T) {
	inputs := []string{
		"(p)",
		"(a, e, i, o, u, y)",
		"SA",
		"END",
	}
	expected := [][]string{
		{"p"},
		{"a", "e", "i", "o", "u", "y"},
		{"a", "e", "i"},
		{"$"},
	}
	I := Interpreter{
		vars: make(map[string][]string),
	}
	I.vars["ALL"] = []string{"a", "e", "i", "o", "u", "y", "cz", "dz"}
	I.vars["SA"] = []string{"a", "e", "i"}
	I.vars["END"] = []string{"$"}

	for i := 0; i < len(inputs); i++ {
		has := expected[i]
		want, _ := I.asConstraint(inputs[i])
		if ok := reflect.DeepEqual(has, want); !ok {
			t.Errorf("error: want %v; has %v", want, has)
		}
	}
}

// asConstraint throws an error when context is not enclosed in ( ).
func TestAsConstraintError(t *testing.T) {
	cases := []struct {
		name, ctx string
	}{
		{"format-wrong", "(k,g,ng"},
		{"missing-char", "+k+g+ng)"},
	}
	i := Interpreter{
		vars: make(map[string][]string),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := i.asConstraint(c.ctx)
			if err == nil {
				t.Errorf("expected method to fail on input %s", c.ctx)
			}
		})
	}
}

// Check if difference context is established as specified.
func TestAsDifference(t *testing.T) {
	I := Interpreter{
		vars: make(map[string][]string),
	}
	I.vars["ALL"] = []string{"ng", "cz", "dz", "p", "$"}
	I.vars["SA1"] = []string{"cz", "dz"}
	I.vars["SA2"] = []string{"ng"}
	I.vars["END"] = []string{"$"}

	inputs := []string{
		"-(ng)",
		"-(cz, dz)",
		"-SA1-SA2",
		"-END",
	}
	expected := [][]string{
		{"cz", "dz", "p", "$"},
		{"ng", "p", "$"},
		{"p", "$"},
		{"ng", "cz", "dz", "p"},
	}

	for i := 0; i < len(inputs); i++ {
		has := expected[i]
		want, _ := I.asDifference(inputs[i])
		if ok := reflect.DeepEqual(has, want); !ok {
			t.Errorf("error: want %v; has %v", want, has)
		}
	}
}

// asDifference throws an error when context is not enclosed in ( ).
func TestAsDifferenceError(t *testing.T) {
	cases := []struct {
		name, ctx string
	}{
		{"format-wrong", "-(p, t, k"},
		{"missing-char", "-p-t"},
	}
	i := Interpreter{
		vars: make(map[string][]string),
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := i.asDifference(c.ctx)
			if err == nil {
				t.Errorf("expected method to fail on input %s", c.ctx)
			}
		})
	}
}

// Check if right context is returned based on the format of the input.
func TestContext(t *testing.T) {
	inputs := []string{
		"PUSTY",
		"-(cz, dz)",
		"(a, e, i, o, u, y)",
		"-SA-SA2",
		"SB",
	}
	expected := [][]string{
		nil,
		{"a", "e", "i", "o", "u", "y"},
		{"a", "e", "i", "o", "u", "y"},
		{"cz", "dz"},
		{"cz", "dz"},
	}
	I := Interpreter{
		vars: make(map[string][]string),
	}
	I.vars["ALL"] = []string{"a", "e", "i", "o", "u", "y", "cz", "dz"}
	I.vars["PUSTY"] = []string{"*"}
	I.vars["SA"] = []string{"a", "e", "i"}
	I.vars["SA2"] = []string{"o", "u", "y"}
	I.vars["SB"] = []string{"cz", "dz"}

	for i := 0; i < len(inputs); i++ {
		has := expected[i]
		want, _ := I.context(inputs[i])
		if ok := reflect.DeepEqual(has, want); !ok {
			t.Errorf("error: want %v; has %v", want, has)
		}
	}
}

// Error is returned when Interpreter.vars has no "ALL" key.
func TestContextError(t *testing.T) {
	I := Interpreter{
		vars: make(map[string][]string),
	}
	_, err := I.context("(a, b, c)")
	if err == nil {
		t.Errorf("error was not raised when \"ALL\" is not set")
	}
}
