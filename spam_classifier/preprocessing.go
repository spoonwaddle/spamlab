package main

import (
	"regexp"
	"strings"
)

func RegexpTokenizer(wordPattern string) func(string) []string {
	wordRegexp := regexp.MustCompile(wordPattern)
	return func(text string) []string {
		return wordRegexp.FindAllString(text, -1)
	}
}

func WithoutStopwords(words []string, stopwords map[string]bool) []string {
	filtered := []string{}
	for _, word := range words {
		if _, isStopword := stopwords[word]; !isStopword {
			filtered = append(filtered, word)
		}
	}
	return filtered
}

func Preprocessor() func(string) Text {
	tokenize := RegexpTokenizer("[a-zA-z]+")
	removeStopwords := func(words []string) []string {
		return WithoutStopwords(words, STOPWORDS)
	}
	return func(text string) Text {
		text = strings.ToLower(text)
		words := tokenize(text)
		return removeStopwords(words)
	}
}
