package main_test

import (
	"reflect"
	. "task2"
	"testing"
)

func TestCountWords(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]int
	}{
		{
			input:    "Hello world",
			expected: map[string]int{"hello": 1, "world": 1},
		},
		{
			input:    "Hello hello world",
			expected: map[string]int{"hello": 2, "world": 1},
		},
		{
			input:    "Lorem ipsum dolor sit amet",
			expected: map[string]int{"lorem": 1, "ipsum": 1, "dolor": 1, "sit": 1, "amet": 1},
		},
		{
			input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit! Lorem ipsum dolor sit amet",
			expected: map[string]int{"lorem": 2, "ipsum": 2, "dolor": 2, "sit": 2, "amet": 2, "consectetur": 1, "adipiscing": 1, "elit": 1},
		},
		{
			input:    "Hello sir, hello madam, and hello world!",
			expected: map[string]int{"hello": 3, "sir": 1, "madam": 1, "and": 1, "world": 1},
		},
	}

	for _, test := range tests {
		result := CountWords(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Input: %s\nExpected: %v\nGot: %v\n", test.input, test.expected, result)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "abacaba",
			expected: true,
		},
		{
			input:    "abacab",
			expected: false,
		},
		{
			input:    "Hello, world!",
			expected: false,
		},
		{
			input:    "Hello, olleH",
			expected: true,
		},
		{
			input:    "This is not a palindrome",
			expected: false,
		},
		{
			input:    "This is Si sIHt",
			expected: true,
		},
	}

	for _, test := range tests {
		result := IsPalindrome(test.input)
		if result != test.expected {
			t.Errorf("Input: %s\nExpected: %v\nGot: %v\n", test.input, test.expected, result)
		}
	}
}
