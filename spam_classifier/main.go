package main

import (
	"fmt"
)

func main() {
	classifier, _ := SpamClassifierFromRedisUrl("127.0.0.1:6379")
	fmt.Println("Clearing model...")
	classifier.dist.ResetCounts()
	fmt.Println("Done.")
	fmt.Println("Training...")
	corpus := EnronSpamCorpus(1, 2, 3, 4, 5, 6)
	classifier.StreamTrain(corpus)
	fmt.Println("Done.")
	fmt.Printf("Listening...")
	fmt.Println(listen(classifier, ":80"))
	fmt.Println("Done.")
}
