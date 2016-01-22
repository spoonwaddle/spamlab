package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var subcommands = map[string]func(){
	"train":    TrainCommand,
	"reset":    ResetCommand,
	"classify": ClassifyCommand,
	"server":   ServerCommand,
}

func ServerCommand() {
	var redisUrl, serverUrl string
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	serverFlags.StringVar(&redisUrl, "redis", os.Getenv("REDIS_URL"), "URL of Redis instance being used to store model.")
	serverFlags.StringVar(&serverUrl, "addr", "0.0.0.0:8080", "Address to host spam classifier server from.")
	serverFlags.Parse(argsAfterSubcommand())
	classifier := parseClassifier(redisUrl)
	fmt.Println("Listening on", serverUrl)
	err := listen(classifier, serverUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ResetCommand() {
	var redisUrl string
	resetFlags := flag.NewFlagSet("reset", flag.ExitOnError)
	resetFlags.StringVar(&redisUrl, "redis", os.Getenv("REDIS_URL"), "URL of Redis instance being used to store model.")
	resetFlags.Parse(argsAfterSubcommand())
	classifier := parseClassifier(redisUrl)
	err := classifier.Reset()
	if err != nil {
		fmt.Println("Bad things happened during reset:", err)
		os.Exit(1)
	} else {
		fmt.Println("All keys flushed from redis server.")
	}
}

func ClassifyCommand() {
	var redisUrl string
	classifyFlags := flag.NewFlagSet("classify", flag.ExitOnError)
	classifyFlags.StringVar(&redisUrl, "redis", os.Getenv("REDIS_URL"), "URL of Redis instance being used to store model.")
	classifyFlags.Parse(argsAfterSubcommand())
	classifier := parseClassifier(redisUrl)
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	label, err := classifier.Classify(string(bytes))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(label)
}

type corpusIndices []int

func (indices *corpusIndices) String() string {
	var ix []string
	for _, i := range *indices {
		ix = append(ix, strconv.Itoa(i))
	}
	return strings.Join(ix, ",")
}

func (indices *corpusIndices) Set(value string) error {
	if len(*indices) > 0 {
		return errors.New("indices flag set multiple times!")
	}
	for _, indexStr := range strings.Split(value, ",") {
		index, err := strconv.Atoi(indexStr)
		if err != nil {
			return errors.New("indices flag must be comma-separated list of integers from 1 to 6")
		}
		*indices = append(*indices, index)
	}
	return nil
}

func TrainCommand() {
	trainingFlags := flag.NewFlagSet("train", flag.ExitOnError)
	var indices corpusIndices
	var hamGlob string
	var spamGlob string
	var redisUrl string

	trainingFlags.StringVar(&spamGlob, "spam", "", "Glob to match training files labelled spam.")
	trainingFlags.StringVar(&hamGlob, "ham", "", "Glob to match training files labelled ham.")
	trainingFlags.StringVar(&redisUrl, "redis", os.Getenv("REDIS_URL"), "URL of Redis instance being used to store model.")
	trainingFlags.Var(&indices, "enron", "comma-separated list of Enron Spam Dataset indices in range [1-6]")
	trainingFlags.Parse(argsAfterSubcommand())

	classifier := parseClassifier(redisUrl)

	var corpora []chan TrainingSample
	if hamGlob != "" {
		c, err := GlobCorpus(HAM, hamGlob)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		corpora = append(corpora, c)
	}
	if spamGlob != "" {
		c, err := GlobCorpus(SPAM, spamGlob)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		corpora = append(corpora, c)
	}
	if len(indices) > 0 {
		c, err := EnronSpamCorpus(indices...)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		corpora = append(corpora, c)
	}
	if len(corpora) == 0 {
		fmt.Println("Nothing to train on!")
		os.Exit(1)
	}
	corpus := merge(corpora)
	classifier.StreamTrain(corpus)
}

func getSubcommandNames(subcommands map[string]func()) []string {
	keys := make([]string, 0, len(subcommands))
	for k, _ := range subcommands {
		keys = append(keys, k)
	}
	return keys
}

func parseSubcommand() func() {
	names := getSubcommandNames(subcommands)
	if len(os.Args) < 2 {
		fmt.Println("Must supply one of the following subcommands:", names)
		os.Exit(1)
	}
	subcommand, commandFound := subcommands[os.Args[1]]
	if !commandFound {
		fmt.Printf("Invalid subcommand: %s. Must be in %s\n", os.Args[1], names)
		os.Exit(1)
	}
	return subcommand
}

func parseClassifier(redisUrl string) SpamClassifier {
	if redisUrl == "" {
		fmt.Println("Must provide redis url via either `redis` flag or REDIS_URL environment variable!")
		os.Exit(1)
	}
	classifier, err := SpamClassifierFromRedisUrl(redisUrl)
	if err != nil {
		fmt.Println("Could not create SpamClassifier from redis url:", redisUrl)
		os.Exit(1)
	}
	return classifier
}

func argsAfterSubcommand() []string {
	var remainingArgs []string
	if len(os.Args) > 2 {
		remainingArgs = os.Args[2:]
	}
	return remainingArgs
}

func main() {
	parseSubcommand()()
}
