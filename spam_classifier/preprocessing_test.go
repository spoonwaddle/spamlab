package main

import (
	"reflect"
	"testing"
)

func TestRegexpTokenizer(t *testing.T) {
	tokenizer := RegexpTokenizer("[A-Z]+")
	expected, result := tokenizer("OH hello THERE"), []string{"OH", "THERE"}
	if !reflect.DeepEqual(expected, result) {
		t.Error("Expected", expected, "got", result)
	}
}

func TestWithoutStopwords(t *testing.T) {
	words := []string{"this", "is", "not", "right"}
	stopwords := map[string]bool{"not": true}
	expected := []string{"this", "is", "right"}
	result := WithoutStopwords(words, stopwords)
	if !reflect.DeepEqual(expected, result) {
		t.Error("Expected", expected, "got", result)
	}
}

func TestDefaultPreprocessor(t *testing.T) {
	preprocess := Preprocessor()
	text := "LeT Us GO        2 THE PARK"
	var expected Text = Text{"let", "park"}
	result := preprocess(text)
	if !reflect.DeepEqual(expected, result) {
		t.Error("Expected", expected, "got", result)
	}
}
