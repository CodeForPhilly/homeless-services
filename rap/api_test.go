package rap

import "testing"

func TestGetResources(t *testing.T) {
}

func TestFilterParser(t *testing.T) {
}

func TestLogicalOperatorConverter(t *testing.T) {
	for _, tt := range locTests {
		actual, err := logicalOperatorConverter(tt.eo)
		if err != tt.err || actual != tt.expected {
			t.Errorf("Logical operator converter(%s): expected %s, actual %s", tt.eo, tt.expected, actual)
		}
	}
}

type locTest struct {
	eo       string
	expected string
	err      error
}

var locTests = []locTest{
	{"gt", ">", nil},
	{"lt", "<", nil},
	{"eq", "=", nil},
	{"le", "<=", nil},
	{"ge", ">=", nil},
	{"ne", "", ErrUnknownOperator}, //not sure why this is failing
}
