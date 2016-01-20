package main

import (
	"fmt"
	"net/http"
)

func classifyEndpoint(classifier SpamClassifier) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		text := request.URL.Query().Get("text")
		if len(text) == 0 {
			http.Error(
				response,
				"Must provide `text` paramter containing text to classify.",
				400)
			return
		}
		label, err := classifier.Classify(text)
		if err != nil {
			http.Error(
				response,
				err.Error(),
				500)
			return
		}
		fmt.Fprintf(response, label.String())
	}
}

func listen(classifier SpamClassifier, listenAddr string) error {
	http.HandleFunc("/classify", classifyEndpoint(classifier))
	return http.ListenAndServe(listenAddr, nil)
}
