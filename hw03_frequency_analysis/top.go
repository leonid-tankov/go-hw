package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regExpr = regexp.MustCompile(`[-a-zA-Zа-яА-Я]+`)

type wordCount struct {
	word  string
	count int
}

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
	wordCountSlice := make([]wordCount, len(wordsCountMap))
	var counter int
	for word, count := range wordsCountMap {
		wordCountSlice[counter] = wordCount{word, count}
		counter++
	}
	sort.Slice(wordCountSlice, func(i, j int) bool {
		if wordCountSlice[i].count > wordCountSlice[j].count {
			return true
		}
		if wordCountSlice[i].count == wordCountSlice[j].count && wordCountSlice[i].word < wordCountSlice[j].word {
			return true
		}
		return false
	})
	var length int
	if len(wordCountSlice) >= 10 {
		length = 10
	} else {
		length = len(wordCountSlice)
	}
	result := make([]string, length)
	for i := 0; i < length; i++ {
		result[i] = wordCountSlice[i].word
	}

	return result
}
