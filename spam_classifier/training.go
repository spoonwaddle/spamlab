package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func addToCorpus(glob string, label SpamLabel, corpus chan TrainingSample) {
	files, err := filepath.Glob(glob)
	if err != nil {
		panic(fmt.Sprintf("%s: %s", "Could not generate paths from glob", err))
	}
	for _, f := range files {
		text, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Printf("Skipping %s due to error: %s\n", f, err)
		} else {
			corpus <- TrainingSample{label: label, document: string(text)}
		}
	}
}

func addDocumentToCorpus(doc string, label SpamLabel, corpus chan TrainingSample) {
	corpus <- TrainingSample{label: label, document: doc}
}

func TrainFromFiles(classifier SpamClassifier, spamGlob string, hamGlob string) {
	corpus := make(chan TrainingSample)
	classifier.StreamTrain(corpus)
	addToCorpus(spamGlob, SPAM, corpus)
	addToCorpus(hamGlob, HAM, corpus)
}
