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

const ENRON_DATASET_BASE_URL = "http://www.aueb.gr/users/ion/data/enron-spam/preprocessed/enron"

var LABEL2NAMESPACE = map[SpamLabel]string{
	SPAM: "spam:",
	HAM:  "ham:",
}
