// Package threshold parses and evaluates pass/fail conditions on a run's
// metrics, so http-runner can gate CI pipelines (exit non-zero when a latency
// or success-rate budget is violated).
//
// A condition is written as "<metric><op><value>", and several are joined with
// commas, e.g. "p99>500ms,success<99,errors>0". Each condition describes a
// FAILURE: the run fails if any condition holds.
package threshold

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Kind describes how a metric's value is written and formatted.
type Kind int

const (
	KindDuration Kind = iota // written as a Go duration (e.g. "500ms"); compared in seconds
	KindPercent              // a percentage 0-100 (e.g. success rate)
	KindFloat                // a plain float (e.g. requests/sec)
	KindInt                  // a whole number (e.g. error count)
)

// metrics maps a condition metric name to its kind. Duration metrics are
// compared in seconds, matching the report fields.
var metrics = map[string]Kind{
	"p50":     KindDuration,
	"p90":     KindDuration,
	"p95":     KindDuration,
	"p99":     KindDuration,
	"avg":     KindDuration,
	"min":     KindDuration,
	"max":     KindDuration,
	"ttfb":    KindDuration,
	"success": KindPercent,
	"rps":     KindFloat,
	"errors":  KindInt,
}

// ops are the supported comparison operators, longest first so ">=" is matched
// before ">".
var ops = []string{">=", "<=", "==", "!=", ">", "<"}

// Condition is a single parsed failure threshold.
type Condition struct {
	Metric string  // metric name (e.g. "p99")
	Op     string  // comparison operator
	Value  float64 // threshold value; durations normalised to seconds
	Kind   Kind    // kind of the metric
	Raw    string  // original text, for error messages
}

// Parse parses a comma-separated threshold spec into conditions. An empty spec
// yields no conditions.
func Parse(spec string) ([]Condition, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, nil
	}
	var conds []Condition
	for _, tok := range strings.Split(spec, ",") {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		c, err := parseCondition(tok)
		if err != nil {
			return nil, err
		}
		conds = append(conds, c)
	}
	return conds, nil
}

func parseCondition(tok string) (Condition, error) {
	for _, op := range ops {
		i := strings.Index(tok, op)
		if i <= 0 {
			continue
		}
		metric := strings.TrimSpace(tok[:i])
		valStr := strings.TrimSpace(tok[i+len(op):])
		kind, ok := metrics[metric]
		if !ok {
			return Condition{}, fmt.Errorf("unknown metric %q in %q", metric, tok)
		}
		val, err := parseValue(kind, valStr)
		if err != nil {
			return Condition{}, fmt.Errorf("invalid value in %q: %w", tok, err)
		}
		return Condition{Metric: metric, Op: op, Value: val, Kind: kind, Raw: tok}, nil
	}
	return Condition{}, fmt.Errorf("no operator (one of >, <, >=, <=, ==, !=) in %q", tok)
}

func parseValue(kind Kind, s string) (float64, error) {
	if kind == KindDuration {
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}
		return d.Seconds(), nil
	}
	return strconv.ParseFloat(s, 64)
}

// Evaluate returns a message for every condition that holds against values
// (metric name -> actual value; durations in seconds). An empty result means
// all thresholds passed. Conditions on metrics absent from values are skipped.
func Evaluate(conds []Condition, values map[string]float64) []string {
	var fails []string
	for _, c := range conds {
		actual, ok := values[c.Metric]
		if !ok {
			continue
		}
		if compare(actual, c.Op, c.Value) {
			fails = append(fails, fmt.Sprintf("%s (actual %s)", c.Raw, formatActual(c.Kind, actual)))
		}
	}
	return fails
}

func compare(a float64, op string, b float64) bool {
	switch op {
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	case "==":
		return a == b
	case "!=":
		return a != b
	}
	return false
}

func formatActual(kind Kind, v float64) string {
	switch kind {
	case KindDuration:
		return fmt.Sprintf("%.6fs", v)
	case KindPercent:
		return fmt.Sprintf("%.2f%%", v)
	case KindInt:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%.2f", v)
	}
}
