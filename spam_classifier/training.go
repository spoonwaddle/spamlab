package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func addDocumentToCorpus(doc string, label SpamLabel, corpus chan TrainingSample) {
	corpus <- TrainingSample{label: label, document: doc}
}

func GlobCorpus(label SpamLabel, glob string) (chan TrainingSample, error) {
	corpus := make(chan TrainingSample)
	files, err := filepath.Glob(glob)
	if err != nil {
		msg := fmt.Sprintf("Could not generate paths from glob (%s): %s", glob, err)
		return nil, errors.New(msg)
	}
	go func() {
		for _, f := range files {
			text, err := ioutil.ReadFile(f)
			if err != nil {
				fmt.Printf("Skipping %s due to error: %s\n", f, err)
			} else {
				corpus <- TrainingSample{label: label, document: string(text)}
			}
		}
		close(corpus)
	}()
	return corpus, nil
}
