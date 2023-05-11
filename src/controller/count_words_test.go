package controller

import (
	"reflect"
	"testing"
)

func TestCountWords(t *testing.T) {
	input := "Use.proper tool for proper|thing. Use.proper tool for proper|thing."
	expected := map[string]int{
		"Use":    2,
		"proper": 4,
		"tool":   2,
		"for":    2,
		"thing":  2,
	}

	result := countWords(input)
	// HINT: reflect package to have things quickly done
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected word count. Got: %v, Expected: %v", result, expected)
	}
}
