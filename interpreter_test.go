package prg2p

import (
	"os"
	"reflect"
	"testing"
)

// Check if the interpreter handles the file with rules without raising an
// error.
func TestParser(t *testing.T) {
	i := &Interpreter{
		vars: make(map[string][]string),
	}
	f, _ := os.Open("../vendor/g2p-rules.txt")
	defer f.Close()
	err := i.Scan(f)
	if err != nil {
		t.Errorf("%v", err)
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
	I := Interpreter{
		vars: make(map[string][]string),
	}
	ctx := "(k,g,ng"
	_, err := I.asConstraint(ctx)
	if err == nil {
		t.Errorf("wrong constraint format hasn't raised an error, %s", ctx)
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
	I := Interpreter{
		vars: make(map[string][]string),
	}
	ctx := "-(p, t, k"
	_, err := I.asDifference(ctx)
	if err == nil {
		t.Errorf("wrong constraint format hasn't raised an error, %s", ctx)
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
