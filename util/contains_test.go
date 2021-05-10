package util

import "testing"

func TestContainsString(t *testing.T) {
	stack := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
	}

	result := ContainsString(stack, "one")
	if result != true {
		t.Errorf("Sum was incorrect, got: %v, want: %v.", result, true)
	}

	result = ContainsString(stack, "four")
	if result != true {
		t.Errorf("Sum was incorrect, got: %v, want: %v.", result, true)
	}

	result = ContainsString(stack, "foo")
	if result != false {
		t.Errorf("Sum was incorrect, got: %v, want: %v.", result, false)
	}
}

func TestContainsOneOfStrings(t *testing.T) {
	stack := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
	}

	result := ContainsOneOfStrings(stack, []string{"one", "two", "three"})
	if result != true {
		t.Errorf("Sum was incorrect, got: %v, want: %v.", result, true)
	}

	result = ContainsOneOfStrings(stack, []string{"foo", "five", "bar"})
	if result != true {
		t.Errorf("Sum was incorrect, got: %v, want: %v.", result, true)
	}

	result = ContainsOneOfStrings(stack, []string{"foo", "bar", "fizz"})
	if result != false {
		t.Errorf("Sum was incorrect, got: %v, want: %v.", result, false)
	}

}
