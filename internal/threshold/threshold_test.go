package threshold

import "testing"

func TestParse_Valid(t *testing.T) {
	conds, err := Parse("p99>500ms, success<99 , errors>0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(conds) != 3 {
		t.Fatalf("expected 3 conditions, got %d", len(conds))
	}
	if conds[0].Metric != "p99" || conds[0].Op != ">" || conds[0].Value != 0.5 {
		t.Errorf("p99 parsed wrong: %+v", conds[0])
	}
	if conds[1].Metric != "success" || conds[1].Op != "<" || conds[1].Value != 99 {
		t.Errorf("success parsed wrong: %+v", conds[1])
	}
}

func TestParse_Empty(t *testing.T) {
	conds, err := Parse("   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conds != nil {
		t.Errorf("expected nil conditions, got %+v", conds)
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []string{
		"p99!500ms",      // no valid operator
		"bogus>1",        // unknown metric
		"p99>notadur",    // bad duration
		"success>notnum", // bad float
		">500ms",         // no metric
	}
	for _, spec := range cases {
		if _, err := Parse(spec); err == nil {
			t.Errorf("expected error for %q, got nil", spec)
		}
	}
}

func TestParse_GreaterEqualBeforeGreater(t *testing.T) {
	conds, err := Parse("p95>=100ms")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conds[0].Op != ">=" {
		t.Errorf("expected op >=, got %q", conds[0].Op)
	}
}

func TestEvaluate(t *testing.T) {
	conds, err := Parse("p99>500ms,success<99,errors>0,rps<10")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	values := map[string]float64{
		"p99":     0.62, // 620ms, violates p99>500ms
		"success": 100,  // ok
		"errors":  0,    // ok
		"rps":     42,   // ok
	}
	fails := Evaluate(conds, values)
	if len(fails) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(fails), fails)
	}
}

func TestEvaluate_AllPass(t *testing.T) {
	conds, _ := Parse("p99>1s,success<50")
	values := map[string]float64{"p99": 0.1, "success": 100}
	if fails := Evaluate(conds, values); len(fails) != 0 {
		t.Errorf("expected no violations, got %v", fails)
	}
}
