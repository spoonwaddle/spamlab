package main

type Text []string
type SpamLabel bool

const SPAM SpamLabel = true
const HAM SpamLabel = false

func (label SpamLabel) String() string {
	if label == SPAM {
		return "SPAM"
	} else {
		return "HAM"
	}
}

const (
	TEXT_COUNT_SUFFIX = "#TEXTS"
	WORD_COUNT_SUFFIX = "#WORDS"
)

var LABEL2NAMESPACE = map[SpamLabel]string{
	SPAM: "spam:",
	HAM:  "ham:",
}
