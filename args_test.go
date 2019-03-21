package cobra

import (
	"strings"
	"testing"
)

type argsTestcase struct {
	exerr  string         // Expected error key (see map[string][string])
	args   PositionalArgs // Args validator
	wValid bool           // Define `ValidArgs` in the command
	rargs  []string       // Runtime args
}

var errStrings = map[string]string{
	"invalid":    `invalid argument "a" for "c"`,
	"unknown":    `unknown command "one" for "c"`,
	"less":       "requires at least 2 arg(s), only received 1",
	"more":       "accepts at most 2 arg(s), received 3",
	"notexact":   "accepts 2 arg(s), received 3",
	"notinrange": "accepts between 2 and 4 arg(s), received 1",
}

func (tc *argsTestcase) test(t *testing.T) {
	c := &Command{
		Use:  "c",
		Args: tc.args,
		Run:  emptyRun,
	}
	if tc.wValid {
		c.ValidArgs = []string{"one", "two", "three"}
	}

	o, e := executeCommand(c, tc.rargs...)

	if len(tc.exerr) > 0 {
		// Expect error
		if e == nil {
			t.Fatal("Expected an error")
		}
		expected, ok := errStrings[tc.exerr]
		if !ok {
			t.Errorf(`key "%s" is not found in map "errStrings"`, tc.exerr)
			return
		}
		if got := e.Error(); got != expected {
			t.Errorf("Expected: %q, got: %q", expected, got)
		}
	} else {
		// Expect success
		if o != "" {
			t.Errorf("Unexpected output: %v", o)
		}
		if e != nil {
			t.Fatalf("Unexpected error: %v", e)
		}
	}
}

func TestArgs(t *testing.T) {
	tests := map[string]argsTestcase{
		"No        |       |  ":                   {"", NoArgs, false, []string{}},
		"No        |       | Arb":                 {"unknown", NoArgs, false, []string{"one"}},
		"No        | Valid | Valid":               {"unknown", NoArgs, true, []string{"one"}},
		"Arbitrary |       | Arb":                 {"", ArbitraryArgs, false, []string{"a", "b"}},
		"Arbitrary | Valid | Valid":               {"", ArbitraryArgs, true, []string{"one", "two"}},
		"Arbitrary | Valid | Invalid":             {"invalid", ArbitraryArgs, true, []string{"a"}},
		"MinimumN  |       | Arb":                 {"", MinimumNArgs(2), false, []string{"a", "b", "c"}},
		"MinimumN  | Valid | Valid":               {"", MinimumNArgs(2), true, []string{"one", "three"}},
		"MinimumN  | Valid | Invalid":             {"invalid", MinimumNArgs(2), true, []string{"a", "b"}},
		"MinimumN  |       | Less":                {"less", MinimumNArgs(2), false, []string{"a"}},
		"MinimumN  | Valid | Less":                {"less", MinimumNArgs(2), true, []string{"one"}},
		"MinimumN  | Valid | LessInvalid":         {"invalid", MinimumNArgs(2), true, []string{"a"}},
		"MaximumN  |       | Arb":                 {"", MaximumNArgs(3), false, []string{"a", "b"}},
		"MaximumN  | Valid | Valid":               {"", MaximumNArgs(2), true, []string{"one", "three"}},
		"MaximumN  | Valid | Invalid":             {"invalid", MaximumNArgs(2), true, []string{"a", "b"}},
		"MaximumN  |       | More":                {"more", MaximumNArgs(2), false, []string{"a", "b", "c"}},
		"MaximumN  | Valid | More":                {"more", MaximumNArgs(2), true, []string{"one", "three", "two"}},
		"MaximumN  | Valid | MoreInvalid":         {"invalid", MaximumNArgs(2), true, []string{"a", "b", "c"}},
		"Exact     |       | Arb":                 {"", ExactArgs(3), false, []string{"a", "b", "c"}},
		"Exact     | Valid | Valid":               {"", ExactArgs(3), true, []string{"three", "one", "two"}},
		"Exact     | Valid | Invalid":             {"invalid", ExactArgs(3), true, []string{"three", "a", "two"}},
		"Exact     |       | InvalidCount":        {"notexact", ExactArgs(2), false, []string{"a", "b", "c"}},
		"Exact     | Valid | InvalidCount":        {"notexact", ExactArgs(2), true, []string{"three", "one", "two"}},
		"Exact     | Valid | InvalidCountInvalid": {"invalid", ExactArgs(2), true, []string{"three", "a", "two"}},
		"Range     |       | Arb":                 {"", RangeArgs(2, 4), false, []string{"a", "b", "c"}},
		"Range     | Valid | Valid":               {"", RangeArgs(2, 4), true, []string{"three", "one", "two"}},
		"Range     | Valid | Invalid":             {"invalid", RangeArgs(2, 4), true, []string{"three", "a", "two"}},
		"Range     |       | InvalidCount":        {"notinrange", RangeArgs(2, 4), false, []string{"a"}},
		"Range     | Valid | InvalidCount":        {"notinrange", RangeArgs(2, 4), true, []string{"two"}},
		"Range     | Valid | InvalidCountInvalid": {"invalid", RangeArgs(2, 4), true, []string{"a"}},
		//DEPRECATED
		"DEPRECATED OnlyValid  | Valid | Valid":        {"", OnlyValidArgs, true, []string{"one", "two"}},
		"DEPRECATED OnlyValid  | Valid | Invalid":      {"invalid", OnlyValidArgs, true, []string{"a"}},
		"DEPRECATED ExactValid | Valid | Valid":        {"", ExactValidArgs(3), true, []string{"two", "three", "one"}},
		"DEPRECATED ExactValid | Valid | InvalidCount": {"notexact", ExactValidArgs(2), true, []string{"two", "three", "one"}},
		"DEPRECATED ExactValid | Valid | Invalid":      {"invalid", ExactValidArgs(2), true, []string{"two", "a"}},
	}

	for name, tc := range tests {
		t.Run(name, tc.test)
	}
}

// Takes(No)Args

func TestRootTakesNoArgs(t *testing.T) {
	rootCmd := &Command{Use: "root", Run: emptyRun}
	childCmd := &Command{Use: "child", Run: emptyRun}
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "illegal", "args")
	if err == nil {
		t.Fatal("Expected an error")
	}

	got := err.Error()
	expected := `unknown command "illegal" for "root"`
	if !strings.Contains(got, expected) {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestRootTakesArgs(t *testing.T) {
	rootCmd := &Command{Use: "root", Args: ArbitraryArgs, Run: emptyRun}
	childCmd := &Command{Use: "child", Run: emptyRun}
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "legal", "args")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestChildTakesNoArgs(t *testing.T) {
	rootCmd := &Command{Use: "root", Run: emptyRun}
	childCmd := &Command{Use: "child", Args: NoArgs, Run: emptyRun}
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "child", "illegal", "args")
	if err == nil {
		t.Fatal("Expected an error")
	}

	got := err.Error()
	expected := `unknown command "illegal" for "root child"`
	if !strings.Contains(got, expected) {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestChildTakesArgs(t *testing.T) {
	rootCmd := &Command{Use: "root", Run: emptyRun}
	childCmd := &Command{Use: "child", Args: ArbitraryArgs, Run: emptyRun}
	rootCmd.AddCommand(childCmd)

	_, err := executeCommand(rootCmd, "child", "legal", "args")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
