package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

func buildEnronTarballUrl(tarballIndex int) (string, error) {
	if (tarballIndex > 6) && (tarballIndex < 0) {
		return "", errors.New("Enron dataset file index must be in range [1-6].")
	}
	return fmt.Sprintf("%s%d.tar.gz", ENRON_DATASET_BASE_URL, tarballIndex), nil
}

func readTarballFromUrl(tarballUrl string) (*tar.Reader, error) {
	response, err := http.Get(tarballUrl)
	if err != nil {
		downloadError := fmt.Sprintf("Unable to download %s : %s", tarballUrl, err)
		return nil, errors.New(downloadError)
	}
	gzFile, err := gzip.NewReader(response.Body)
	if err != nil {
		gzError := fmt.Sprintf("Unable to create Gzip reader for %s: %s", tarballUrl, err)
		return nil, errors.New(gzError)
	}
	return tar.NewReader(gzFile), nil
}

func readerToString(reader *tar.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}

func corpusFromTar(tarReader *tar.Reader) chan TrainingSample {
	corpus := make(chan TrainingSample)
	go func() {
		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err == nil && header.Typeflag == tar.TypeReg {
				if strings.HasSuffix(header.Name, ".spam.txt") {
					text := readerToString(tarReader)
					fmt.Println("Training on", header.Name)
					addDocumentToCorpus(text, SPAM, corpus)
				} else if strings.HasSuffix(header.Name, ".ham.txt") {
					text := readerToString(tarReader)
					fmt.Println("Training on", header.Name)
					addDocumentToCorpus(text, HAM, corpus)
				}
			} else {
				continue
			}
		}
		close(corpus)
	}()
	return corpus
}

func merge(corpora []chan TrainingSample) chan TrainingSample {
	var wg sync.WaitGroup
	out := make(chan TrainingSample)
	output := func(c <-chan TrainingSample) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(corpora))
	for _, c := range corpora {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func EnronSpamCorpus(indices ...int) (chan TrainingSample, error) {
	var corpora []chan TrainingSample
	for _, i := range indices {
		url, err := buildEnronTarballUrl(i)
		if err != nil {
			return nil, err
		}
		tarball, err := readTarballFromUrl(url)
		if err != nil {
			return nil, err
		}
		corpora = append(corpora, corpusFromTar(tarball))
	}
	return merge(corpora), nil
}
