package main

type TrainingSample struct {
	document string
	label    SpamLabel
}

type SpamClassifier struct {
	dist     WordProbabilityDist
	tokenize func(string) Text
}

func SpamClassifierFromRedisUrl(redisUrl string) (SpamClassifier, error) {
	dist, err := NewRedisWordProbabilityDist(redisUrl)
	if err != nil {
		return SpamClassifier{}, err
	}
	tokenize := Preprocessor()
	return SpamClassifier{
		dist:     dist,
		tokenize: tokenize,
	}, nil
}

func (classifier SpamClassifier) Train(samples ...TrainingSample) error {
	for _, sample := range samples {
		words := classifier.tokenize(sample.document)
		err := classifier.dist.Update(words, sample.label)
		if err != nil {
			return err
		}
	}
	return nil
}

func (classifier SpamClassifier) Classify(document string) (SpamLabel, error) {
	words := classifier.tokenize(document)
	label, err := classifier.dist.MostLikelyLabel(words)
	if err != nil {
		return false, err
	}
	return label, nil
}

func (classifier SpamClassifier) StreamTrain(samples chan TrainingSample) {
	for sample := range samples {
		classifier.Train(sample)
	}
}

func (classifier SpamClassifier) Reset() error {
	return classifier.dist.ResetCounts()
}
