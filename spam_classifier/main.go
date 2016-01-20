package main

import (
	"fmt"
)

func main() {
	classifier, _ := SpamClassifierFromRedisUrl("127.0.0.1:6379")
	fmt.Printf("Clearing model...")
	classifier.dist.ResetCounts()
	fmt.Println("Done.")
	fmt.Printf("Training...")
	TrainFromFiles(
		classifier,
		"/home/vagrant/training_data/spam/*",
		"/home/vagrant/training_data/ham/*",
	)
	fmt.Println("Done.")
	fmt.Printf("Listening...")
	fmt.Println(listen(classifier, ":80"))
	fmt.Println("Done.")
}
