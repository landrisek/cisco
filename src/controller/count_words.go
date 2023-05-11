package controller

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"unicode"
)

func countWords(input string) map[string]int {
	words := strings.FieldsFunc(input, func(r rune) bool {
		// HINT: I was trying to avoid regex due their requirement on performance
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	count := make(map[string]int)
	// HINT: dirty, not good readable way to count all trues from anonymous function
	for _, word := range words {
		count[word]++
	}

	return count
}

// Sort the word counts by count in descending order
type toSort struct {
	word  string
	count int
}

func CountWords(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// HINT: normalization
	normalized := strings.ToLower(string(file))

	counted := countWords(normalized)

	sorted := make([]toSort, 0, len(counted))
	for word, count := range counted {
		sorted = append(sorted, toSort{word, count})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].count > sorted[j].count
	})

	for _, item := range sorted {
		fmt.Printf("%d %s\n", item.count, item.word)
	}
	return nil
}
