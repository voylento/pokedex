package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "  ",
			expected: []string{},
		},
		{
			input: "",
			expected: []string{},
		},  
		{
			input: "hello     world\n testity test test",
			expected: []string{"hello", "world", "testity", "test", "test"},
		},  
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("FAIL:\nINPUT: %v\nEXPECTED length: %v\nACTUAL length: %v", c.input, len(c.expected), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]
			if word != expected {
				t.Errorf("FAIL:\nINPUT: %v\nEXPECTED: %v\nACTUAL: %v", c.input, expected, word)
			}
		}
	}
}
