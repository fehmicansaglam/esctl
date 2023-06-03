package output

import (
	"sort"
	"testing"
)

func TestParseDataSize(t *testing.T) {
	testCases := []struct {
		name     string
		sizeStr  string
		expected float64
	}{
		{"Bytes", "10b", 10},
		{"Kilobytes", "10kb", 10 * 1024},
		{"Megabytes", "10mb", 10 * 1024 * 1024},
		{"Gigabytes", "10gb", 10 * 1024 * 1024 * 1024},
		{"Terabytes", "10tb", 10 * 1024 * 1024 * 1024 * 1024},
		{"Fractional kilobytes", "5.6kb", 5.6 * 1024},
		{"Invalid unit", "10ab", 0},
		{"Invalid value", "ab10", 0},
		{"Mixed case", "10Kb", 10 * 1024},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseDataSize(tc.sizeStr)
			if err != nil && tc.expected != 0 {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("parseDataSize(%s) = %f, want %f", tc.sizeStr, result, tc.expected)
			}
		})
	}
}

func TestSortDataSize(t *testing.T) {
	testCases := []struct {
		name     string
		inputs   []string
		expected []string
	}{
		{
			"Sort data sizes",
			[]string{"10gb", "5kb", "3mb", "2tb", "1b"},
			[]string{"1b", "5kb", "3mb", "10gb", "2tb"},
		},
		{
			"Sort with same sizes different units",
			[]string{"1024b", "1kb"},
			[]string{"1024b", "1kb"},
		},
		{
			"Sort with fractional sizes",
			[]string{"1.5kb", "1500b"},
			[]string{"1500b", "1.5kb"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.SliceStable(tc.inputs, func(i, j int) bool {
				return sortDataSize(tc.inputs[i], tc.inputs[j])
			})
			for i, v := range tc.inputs {
				if v != tc.expected[i] {
					t.Errorf("unexpected sort result: got %v, want %v", tc.inputs, tc.expected)
					break
				}
			}
		})
	}
}
