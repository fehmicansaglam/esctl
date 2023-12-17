package es

import "testing"

func TestExtractFieldAndValue(t *testing.T) {
	tests := []struct {
		input         string
		expectedField string
		expectedValue string
		expectError   bool
	}{
		{"field:value", "field", "value", false},
		{"field:with:colons", "field", "with:colons", false},
		{"invalidformat", "", "", true},
	}

	for _, test := range tests {
		field, value, err := extractFieldAndValue(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %s, but got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %s: %v", test.input, err)
			}

			if field != test.expectedField || value != test.expectedValue {
				t.Errorf("For input %s, expected (%s, %s), but got (%s, %s)", test.input, test.expectedField, test.expectedValue, field, value)
			}
		}
	}
}
