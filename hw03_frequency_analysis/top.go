package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regExpr = regexp.MustCompile(`[-a-zA-Zа-яА-Я]+`)

func Top10(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	words := strings.Fields(text)
	wordsCountMap := map[string]int{}
	for _, word := range words {
		if word == "-" {
			continue
		}
		resultWord := regExpr.FindString(strings.ToLower(word))
		wordsCountMap[resultWord]++
	}
	wordSlice := make([]string, len(wordsCountMap))
	var counter int
	for word := range wordsCountMap {
		wordSlice[counter] = word
		counter++
	}
	sort.Slice(wordSlice, func(i, j int) bool {
		if wordsCountMap[wordSlice[i]] > wordsCountMap[wordSlice[j]] {
			return true
		}
		if wordsCountMap[wordSlice[i]] == wordsCountMap[wordSlice[j]] && wordSlice[i] < wordSlice[j] {
			return true
		}
		return false
	})
	if len(wordSlice) < 10 {
		return wordSlice
	}
	result := make([]string, 10)
	for i := 0; i < 10; i++ {
		result[i] = wordSlice[i]
	}

	return result
}
