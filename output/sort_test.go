package output

import (
	"sort"
	"testing"
)

func TestParseDataSize(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
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
			result, err := parseDataSize(tc.input)
			if err != nil && tc.expected != 0 {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("parseDataSize(%s) = %f, want %f", tc.input, result, tc.expected)
			}
		})
	}
}

type TestCase struct {
	name     string
	input    []string
	expected []string
}

func testSort(t *testing.T, testCases []TestCase, sortFunc func(a, b string) bool) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.SliceStable(tc.input, func(i, j int) bool {
				return sortFunc(tc.input[i], tc.input[j])
			})

			for i, v := range tc.input {
				if v != tc.expected[i] {
					t.Errorf("unexpected sort result: got %v, want %v", tc.input, tc.expected)
					break
				}
			}
		})
	}
}

func TestSortText(t *testing.T) {
	testCases := []TestCase{
		{
			"Sort texts in natural order",
			[]string{
				"cluster-10-node-1",
				"cluster-1-node-2",
				"cluster-1-node-10",
				"cluster-3",
				"cluster-2-node-5",
				"cluster-2-node-10",
				"cluster-1",
				"cluster-1-node-3",
				"cluster-2-node-1",
			},
			[]string{
				"cluster-1",
				"cluster-1-node-2",
				"cluster-1-node-3",
				"cluster-1-node-10",
				"cluster-2-node-1",
				"cluster-2-node-5",
				"cluster-2-node-10",
				"cluster-3",
				"cluster-10-node-1",
			},
		},
	}

	testSort(t, testCases, sortText)
}

func TestSortDataSize(t *testing.T) {
	testCases := []TestCase{
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
		{
			"Sort with empty values",
			[]string{"", "10gb", "", "5kb", "", "3mb"},
			[]string{"", "", "", "5kb", "3mb", "10gb"},
		},
	}

	testSort(t, testCases, sortDataSize)
}

func TestSortPercent(t *testing.T) {
	testCases := []TestCase{
		{
			"Ascending",
			[]string{"12%", "33.9%", "99.1%"},
			[]string{"99.1%", "33.9%", "12%"},
		},
		{
			"Descending",
			[]string{"99.1%", "50%", "12%"},
			[]string{"99.1%", "50%", "12%"},
		},
		{
			"Mixed",
			[]string{"0%", "0.5%", "0.2%"},
			[]string{"0.5%", "0.2%", "0%"},
		},
	}

	testSort(t, testCases, sortPercent)
}

func TestSortDate(t *testing.T) {
	testCases := []TestCase{
		{
			"Ascending",
			[]string{"2021-05-23T18:14:29.392Z", "2022-07-12T10:30:45.123Z", "2024-03-01T05:20:15.678Z"},
			[]string{"2021-05-23T18:14:29.392Z", "2022-07-12T10:30:45.123Z", "2024-03-01T05:20:15.678Z"},
		},
		{
			"Descending",
			[]string{"2023-05-23T18:14:29.392Z", "2022-07-12T10:30:45.123Z", "2021-03-01T05:20:15.678Z"},
			[]string{"2021-03-01T05:20:15.678Z", "2022-07-12T10:30:45.123Z", "2023-05-23T18:14:29.392Z"},
		},
		{
			"Mixed",
			[]string{"2024-03-01T05:20:15.678Z", "2022-07-12T10:30:45.123Z", "2023-05-23T18:14:29.392Z"},
			[]string{"2022-07-12T10:30:45.123Z", "2023-05-23T18:14:29.392Z", "2024-03-01T05:20:15.678Z"},
		},
	}

	testSort(t, testCases, sortDate)
}
