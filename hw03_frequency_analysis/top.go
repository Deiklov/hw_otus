package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const TopK = 10

func Top10(input string) []string {
	splitted := strings.Fields(input)
	return topKFrequent(splitted, TopK)
}

func topKFrequent(words []string, k int) []string {
	freq := make(map[string]int, TopK)
	var uniq []string
	for _, elem := range words {
		if _, ok := freq[elem]; !ok {
			uniq = append(uniq, elem)
		}
		freq[elem]++
	}

	sort.Slice(uniq, func(i, j int) bool {
		if freq[uniq[i]] == freq[uniq[j]] {
			return uniq[i] < uniq[j]
		}
		return freq[uniq[i]] > freq[uniq[j]]
	})

	if k > len(uniq) {
		return uniq
	}
	return uniq[:k]
}
