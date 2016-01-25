# SpamLab

## Synopsis

SpamLab is a system for spam classification written in Go. Its core task is to label input text as either "spam" or "ham" (which is apparently generally accepted by academic types who spend time thinking about spam classifiers to be the opposite of "spam"). SpamLab implements a Naive Bayes classifier and uses Redis to store the underlying statistical model with which it draws its conclusions.

SpamLab exposes a command-line interface through the `spam_classifier` binary, which allows you to train and reset your spam classifier model, classify text, or launch a server to classify text via HTTP GET requests.

## Install

The easiest way to start playing with SpamLab is with Vagrant (I'm using v1.8.1). A simple `vagrant up` should install Golang, Redis, and a couple Go dependencies for you and start up the Redis server. Even without Vagrant, setup is pretty simple. You simply clone this repo, install the dependencies (check out bootstrap.sh), and run `go build` in the `spam_classifier/` directory.

## QuickStart

These instructions will take you from cloning the repository, to training your model with a decent amount of initial data to bootstrap the classifier.

```
git clone git@github.com:vroomwaddle/spamlab.git
cd spamlab/
vagrant up
vagrant ssh
```
In the virtual machine...
```
./integration_tests  # make sure everything's running
cd spam_classifier/
go test              # make sure unit tests pass
./spam_classifier train -enron=1,2,3
```

## Tests

To test the system as a whole, you can run `integration_tests.sh`. This script will create a new Redis server, train a classifier on dummy data, classify a couple of documents to make sure it gets the correct results, and then clean up after itself.

To run unit tests, simply run `go test` in the `spam_classifier/` directory.

## Model

SpamLab's classification is based on a Naive Bayes classifier. At its core, it uses word count data to compute whether the text is more likely to be ham or spam based on labelled samples it has seen before. In order to determine the probability of a label given the text, the probability of each word in the text given the label is multiplied together and then multiplied by the probability of the label. The probability of the labels are computed based on the number of each label the classifier has seen in training data. For this reason, it is important to train the model on roughly the same number of spam and ham samples. In the event that a word has never been seen, we pretend we've seen it once to avoid multiplying by zero probabilities ("add-one smoothing").

## Usage

The `spam_classifier` binary has four subcommands: *train*, *classify*, *server*, and *reset*. Each of these commands talks to the underlying Redis instance containing the model, which means they need to know who to talk to. The Redis instance's address can be set either by the `-redis` flag taken by all subcommands or by the `REDIS_URL` environment variable (bootstrap.sh sets this value if you go the Vagrantfile route).

### train

`train` flags:  
-enron=: comma-separated list of Enron Spam Dataset indices in range [1-6]  
-ham="": Glob to match training files labelled ham.  
-redis="127.0.0.1:6379": URL of Redis instance being used to store model.  
-spam="": Glob to match training files labelled spam.  
    
The `-enron` flag here is worth drawing special attention to. The [Enron Spam Dataset](http://www.aueb.gr/users/ion/data/enron-spam/) contains thousands of emails labelled as either spam or ham. These emails are split into six tarballs. The `-enron` flag takes in indices corresponding to which of the tarballs you would like to train from. These tarballs are asyncronously pulled straight from http://www.aueb.gr/users/ion/data/enron-spam/ and streamed directly into the training function without saving the files to disk. There's a small security risk of the files being tampered with, as files are not downloaded over HTTPS or subject to any other integrity checks. However, given this is mostly a learning project with a presumably small audience, the convenience to users outweighs the value to attackers for now.

Example:

`./spam_classifier train -enron=1,2,3,4,5,6`

The `-ham` and `-spam` flags allow you to specify globs corresponding to groups of files containing either spam or ham data depending on the flag. The files themselves will be treated as plaintext for training.

Example:
`./spam_classifier train -spam=/tmp/*.spam -ham=/tmp/*.ham` 

### classify

Classifies a string from stdin. Result written to stdout as either "SPAM" or "HAM".

Example:

```
echo "pharmaceuticals sure are great!" | ./spam_classifier classify
SPAM
```

### reset

Drops all keys from the Redis instance living at `-redis` or `$REDIS_URL`. This effectively resets the model.

Example:

`./spam_classifier reset -redis=127.0.0.1:5555`

### server

`server` flags:  
-addr="0.0.0.0:8080": Address to host spam classifier server from.    
-redis="127.0.0.1:6379": URL of Redis instance being used to store model.   

Starts an HTTP server to classify texts. Queries are sent via a GET request to the `/classify` endpoint. The text to classify is set with the `text` query parameter in the URL.

Example:

`./spam_classifier server -addr=0.0.0.0:8080`

To test the server out, you can `curl` from it in another window...

`curl "0.0.0.0:8080/classify?text=meeting"`
