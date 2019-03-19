package cobra

import (
	"strings"
	"testing"
)

func newCommand(args PositionalArgs, withValid bool) *Command {
	c := &Command{
		Use:  "c",
		Args: args,
		Run:  emptyRun,
	}
	if withValid {
		c.ValidArgs = []string{"one", "two", "three"}
	}
	return c
}

func expectSuccess(output string, err error, t *testing.T) {
	if output != "" {
		t.Errorf("Unexpected output: %v", output)
	}
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func expectError(err error, t *testing.T, ex string) {
	if err == nil {
		t.Fatal("Expected an error")
	}
	expected := map[string]string{
		"valid": `invalid argument "a" for "c"`,
		"no":    `unknown command "one" for "c"`,
		"min":   "requires at least 2 arg(s), only received 1",
		"max":   "accepts at most 2 arg(s), received 3",
		"exact": "accepts 2 arg(s), received 3",
		"range": "accepts between 2 and 4 arg(s), received 1",
	}[ex]
	if got := err.Error(); got != expected {
		t.Errorf("Expected: %q, got: %q", expected, got)
	}
}

// NoArgs

func TestNoArgs(t *testing.T) {
	o, e := executeCommand(newCommand(NoArgs, false))
	expectSuccess(o, e, t)
}

func TestNoArgsWithArgs(t *testing.T) {
	_, e := executeCommand(newCommand(NoArgs, false), "one")
	expectError(e, t, "no")
}

func TestNoArgsWithArgsWithValid(t *testing.T) {
	_, e := executeCommand(newCommand(NoArgs, true), "one")
	expectError(e, t, "no")
}

// ArbitraryArgs

func TestArbitraryArgs(t *testing.T) {
	o, e := executeCommand(newCommand(ArbitraryArgs, false), "a", "b")
	expectSuccess(o, e, t)
}

func TestArbitraryArgsWithValid(t *testing.T) {
	o, e := executeCommand(newCommand(ArbitraryArgs, true), "one", "two")
	expectSuccess(o, e, t)
}

func TestArbitraryArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(ArbitraryArgs, true), "a")
	expectError(e, t, "valid")
}

// MinimumNArgs

func TestMinimumNArgs(t *testing.T) {
	o, e := executeCommand(newCommand(MinimumNArgs(2), false), "a", "b", "c")
	expectSuccess(o, e, t)
}

func TestMinimumNArgsWithValid(t *testing.T) {
	o, e := executeCommand(newCommand(MinimumNArgs(2), true), "one", "three")
	expectSuccess(o, e, t)
}

func TestMinimumNArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(MinimumNArgs(2), true), "a", "b")
	expectError(e, t, "valid")
}

func TestMinimumNArgsWithLessArgs(t *testing.T) {
	_, e := executeCommand(newCommand(MinimumNArgs(2), false), "a")
	expectError(e, t, "min")
}

func TestMinimumNArgsWithLessArgsWithValid(t *testing.T) {
	_, e := executeCommand(newCommand(MinimumNArgs(2), true), "one")
	expectError(e, t, "min")
}

func TestMinimumNArgsWithLessArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(MinimumNArgs(2), true), "a")
	expectError(e, t, "valid")
}

// MaximumNArgs

func TestMaximumNArgs(t *testing.T) {
	o, e := executeCommand(newCommand(MaximumNArgs(3), false), "a", "b")
	expectSuccess(o, e, t)
}

func TestMaximumNArgsWithValid(t *testing.T) {
	o, e := executeCommand(newCommand(MaximumNArgs(2), true), "one", "three")
	expectSuccess(o, e, t)
}

func TestMaximumNArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(MaximumNArgs(2), true), "a", "b")
	expectError(e, t, "valid")
}

func TestMaximumNArgsWithMoreArgs(t *testing.T) {
	_, e := executeCommand(newCommand(MaximumNArgs(2), false), "a", "b", "c")
	expectError(e, t, "max")
}

func TestMaximumNArgsWithMoreArgsWithValid(t *testing.T) {
	_, e := executeCommand(newCommand(MaximumNArgs(2), true), "one", "three", "two")
	expectError(e, t, "max")
}

func TestMaximumNArgsWithMoreArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(MaximumNArgs(2), true), "a", "b", "c")
	expectError(e, t, "valid")
}

// ExactArgs

func TestExactArgs(t *testing.T) {
	o, e := executeCommand(newCommand(ExactArgs(3), false), "a", "b", "c")
	expectSuccess(o, e, t)
}

func TestExactArgsWithValid(t *testing.T) {
	o, e := executeCommand(newCommand(ExactArgs(3), true), "three", "one", "two")
	expectSuccess(o, e, t)
}

func TestExactArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(ExactArgs(3), true), "three", "a", "two")
	expectError(e, t, "valid")
}

func TestExactArgsWithInvalidCount(t *testing.T) {
	_, e := executeCommand(newCommand(ExactArgs(2), false), "a", "b", "c")
	expectError(e, t, "exact")
}

func TestExactArgsWithInvalidCountWithValid(t *testing.T) {
	_, e := executeCommand(newCommand(ExactArgs(2), true), "three", "one", "two")
	expectError(e, t, "exact")
}

func TestExactArgsWithInvalidCountWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(ExactArgs(2), true), "three", "a", "two")
	expectError(e, t, "valid")
}

// RangeArgs

func TestRangeArgs(t *testing.T) {
	o, e := executeCommand(newCommand(RangeArgs(2, 4), false), "a", "b", "c")
	expectSuccess(o, e, t)
}

func TestRangeArgsWithValid(t *testing.T) {
	o, e := executeCommand(newCommand(RangeArgs(2, 4), true), "three", "one", "two")
	expectSuccess(o, e, t)
}

func TestRangeArgsWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(RangeArgs(2, 4), true), "three", "a", "two")
	expectError(e, t, "valid")
}

func TestRangeArgsWithInvalidCount(t *testing.T) {
	_, e := executeCommand(newCommand(RangeArgs(2, 4), false), "a")
	expectError(e, t, "range")
}

func TestRangeArgsWithInvalidCountWithValid(t *testing.T) {
	_, e := executeCommand(newCommand(RangeArgs(2, 4), true), "two")
	expectError(e, t, "range")
}

func TestRangeArgsWithInvalidCountWithValidWithInvalidArgs(t *testing.T) {
	_, e := executeCommand(newCommand(RangeArgs(2, 4), true), "a")
	expectError(e, t, "valid")
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
