package main

import (
	"sort"
	"strings"
)

// Tiered, deterministic match ranking for search results — a reduced port of
// the npm match-sorter algorithm (MIT, Kent C. Dodds) covering only the
// string-only ranking this search uses. Items rank as: equal, then
// starts-with, word-starts-with, contains, acronym, then a closeness score
// between rankEqual and rankEqual+1 for in-order character matches.
func matchSorter(items []map[string]any, value, key string) []map[string]any {
	type ranked struct {
		item  map[string]any
		rank  float64
		value string
	}

	rankedItems := []ranked{}
	for _, item := range items {
		itemValue, _ := item[key].(string)
		rank, rankedValue := matchRank(itemValue, value)
		if rank >= rankMatches {
			rankedItems = append(rankedItems, ranked{item: item, rank: rank, value: rankedValue})
		}
	}

	sort.SliceStable(rankedItems, func(a, b int) bool {
		if rankedItems[a].rank != rankedItems[b].rank {
			return rankedItems[a].rank > rankedItems[b].rank
		}
		return rankedItems[a].value < rankedItems[b].value
	})

	result := make([]map[string]any, 0, len(rankedItems))
	for _, r := range rankedItems {
		result = append(result, r.item)
	}
	return result
}

const (
	rankEqual          = 6.0
	rankStartsWith     = 5.0
	rankWordStartsWith = 4.0
	rankContains       = 3.0
	rankAcronym        = 2.0
	rankMatches        = 1.0
	rankNoMatch        = 0.0
)

// matchRank returns the ranking for how well stringToRank matches testString.
func matchRank(testString, stringToRank string) (float64, string) {
	if len(stringToRank) > len(testString) {
		return rankNoMatch, testString
	}
	if strings.EqualFold(testString, stringToRank) {
		return rankEqual, testString
	}

	testLower := strings.ToLower(testString)
	rankLower := strings.ToLower(stringToRank)

	if strings.HasPrefix(testLower, rankLower) {
		return rankStartsWith, testString
	}
	if strings.Contains(testLower, " "+rankLower) {
		return rankWordStartsWith, testString
	}
	if strings.Contains(testLower, rankLower) {
		return rankContains, testString
	}
	if len(rankLower) == 1 {
		return rankNoMatch, testString
	}
	if strings.Contains(acronym(testLower), rankLower) {
		return rankAcronym, testString
	}
	return closenessRanking(testLower, rankLower), testString
}

// acronym builds the first-letter acronym of a string, splitting on spaces
// and hyphens.
func acronym(s string) string {
	var b strings.Builder
	for _, word := range strings.Split(s, " ") {
		for _, part := range strings.Split(word, "-") {
			if part != "" {
				b.WriteString(part[:1])
			}
		}
	}
	return b.String()
}

// closenessRanking returns a score between rankMatches and rankEqual based on
// how tightly the query characters appear in order in the string: a tight
// cluster of matched characters scores close to rankEqual, a sparse spread
// scores close to rankMatches.
func closenessRanking(testString, stringToRank string) float64 {
	findMatchingCharacter := func(matchChar string, index int) int {
		for j := index; j < len(testString); j++ {
			if testString[j] == matchChar[0] {
				return j + 1
			}
		}
		return -1
	}

	firstIndex := findMatchingCharacter(stringToRank[:1], 0)
	if firstIndex < 0 {
		return rankNoMatch
	}
	charNumber := firstIndex
	for i := 1; i < len(stringToRank); i++ {
		charNumber = findMatchingCharacter(stringToRank[i:i+1], charNumber)
		if charNumber < 0 {
			return rankNoMatch
		}
	}

	spread := charNumber - firstIndex
	if spread == 0 {
		return rankMatches
	}
	return rankMatches + 1.0/float64(spread)
}
