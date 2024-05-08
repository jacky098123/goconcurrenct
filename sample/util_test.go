package sample

import (
	"flag"
	"testing"
)

func Test_BankRound(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		precision int
		expected  float64
	}{
		{
			name:      "round up",
			x:         1.234,
			precision: 2,
			expected:  1.23,
		},
		{
			name:      "round down",
			x:         1.235,
			precision: 2,
			expected:  1.24,
		},
		{
			name:      "round half up",
			x:         1.235,
			precision: 1,
			expected:  1.2,
		},
		{
			name:      "round half down",
			x:         1.245,
			precision: 1,
			expected:  1.2,
		},
		{
			name:      "round half up with 0.5",
			x:         1.25,
			precision: 1,
			expected:  1.2,
		},
		{
			name:      "round half down with 0.5",
			x:         1.35,
			precision: 1,
			expected:  1.4,
		},
		{
			name:      "round half up with 0.5 2",
			x:         1.45,
			precision: 1,
			expected:  1.5,
		},
		{
			name:      "round half down with 0.5",
			x:         1.55,
			precision: 1,
			expected:  1.6,
		},
		{
			name:      "round half up with 0.5",
			x:         1.65,
			precision: 1,
			expected:  1.7,
		},
		{
			name:      "round half down with 0.5",
			x:         1.75,
			precision: 1,
			expected:  1.8,
		},
		{
			name:      "round half up with 0.5",
			x:         1.85,
			precision: 1,
			expected:  1.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BankRound(tt.x, tt.precision)
			if got != tt.expected {
				t.Errorf("BankRound(%f, %d) = %f; want %f", tt.x, tt.precision, got, tt.expected)
			}
		})
	}
}

func Test_vars(t *testing.T) {
	t.Log("Test vars")

	t.Log("appName:", appName)
	t.Log("featureOptions: ", featureOptions)
	t.Log("incremental:", incremental)
	t.Log("args: ", flag.Args())
}
