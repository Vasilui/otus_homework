package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type Item struct {
	World string
	Count int
}

func Top10(input string) []string {
	if input == "" {
		return nil
	}

	words := strings.FieldsFunc(input, Split)
	items := make([]Item, 0)

	m := SliceToMap(words)
	for key, val := range m {
		items = append(items, Item{World: key, Count: val})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].World <= items[j].World
		}
		return items[i].Count > items[j].Count
	})

	res := make([]string, 0, 10)

	for i := range items {
		res = append(res, items[i].World)
		if len(res) == 10 {
			break
		}
	}

	return res
}

func Split(r rune) bool {
	return r == ':' || r == '.' || unicode.IsSpace(r)
}

func SliceToMap(in []string) map[string]int {
	m := make(map[string]int, 0)

	for i := range in {
		word := strings.TrimFunc(strings.ToLower(in[i]), func(r rune) bool {
			return r != '-' && !unicode.IsLetter(r) && !unicode.IsNumber(r)
		})

		if word != "-" && word != "" {
			m[word]++
		}
	}

	return m
}
