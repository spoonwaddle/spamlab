package main

import (
	"github.com/garyburd/redigo/redis"
)

type WordProbabilityDist interface {
	ProbabilityOfLabel(SpamLabel) (float64, error)
	Update(Text, SpamLabel) error
	ProbabilityOfTextGivenLabel(Text, SpamLabel) (float64, error)
	ProbabilityOfTextAndLabel(Text, SpamLabel) (float64, error)
	MostLikelyLabel(text Text) (SpamLabel, error)
	ResetCounts() error
}

type RedisWordProbabilityDist struct {
	label2namespace map[SpamLabel]string
	conn            redis.Conn
	wordCountSuffix string
	textCountSuffix string
}

func NewRedisWordProbabilityDist(redisURL string) (RedisWordProbabilityDist, error) {
	conn, err := redis.Dial("tcp", redisURL)
	if err != nil {
		return RedisWordProbabilityDist{}, err
	} else {
		return RedisWordProbabilityDist{
			conn:            conn,
			label2namespace: LABEL2NAMESPACE,
			wordCountSuffix: WORD_COUNT_SUFFIX,
			textCountSuffix: TEXT_COUNT_SUFFIX,
		}, nil
	}
}

func (dist RedisWordProbabilityDist) countText(prefix string) {
	dist.conn.Send("INCR", dist.textCountSuffix)
	dist.conn.Send("INCR", prefix+dist.textCountSuffix)
}

func (dist RedisWordProbabilityDist) countWords(prefix string, text Text) {
	for _, word := range text {
		dist.conn.Send("INCR", prefix+dist.wordCountSuffix)
		dist.conn.Send("INCR", prefix+word)
	}
}

func (dist RedisWordProbabilityDist) Update(text Text, label SpamLabel) error {
	prefix := dist.label2namespace[label]
	dist.conn.Send("MULTI")
	dist.countText(prefix)
	dist.countWords(prefix, text)
	_, err := dist.conn.Do("EXEC")
	return err
}

func (dist RedisWordProbabilityDist) textToKeys(label SpamLabel, words ...string) []interface{} {
	prefix := dist.label2namespace[label]
	redisKeys := make([]interface{}, len(words))
	for i, word := range words {
		redisKeys[i] = prefix + word
	}
	return redisKeys
}

func (dist RedisWordProbabilityDist) ResetCounts() error {
	_, err := dist.conn.Do("FLUSHALL")
	return err
}

func (dist RedisWordProbabilityDist) ProbabilityOfLabel(label SpamLabel) (float64, error) {
	labelCountKey := dist.label2namespace[label] + dist.textCountSuffix
	totalCountKey := dist.textCountSuffix
	results, err := redis.Ints(dist.conn.Do("MGET", labelCountKey, totalCountKey))
	if err != nil {
		return 0.0, err
	}
	labelCount, totalCount := results[0], results[1]
	return float64(labelCount) / float64(totalCount), nil
}

func (dist RedisWordProbabilityDist) ProbabilityOfTextGivenLabel(text Text, label SpamLabel) (float64, error) {
	labelWordCountKey := dist.label2namespace[label] + dist.wordCountSuffix
	wordKeys := dist.textToKeys(label, text...)
	allKeys := append([]interface{}{labelWordCountKey}, wordKeys...)
	result, err := redis.Ints(dist.conn.Do("MGET", allKeys...))
	if err != nil {
		return 0.0, err
	}
	nWordsInLabel := float64(result[0])
	prob := 1.0
	for _, wordFreq := range result[1:] {
		if wordFreq == 0 {
			wordFreq = 1 // add-one smoothing
		}
		prob *= float64(wordFreq) / nWordsInLabel
	}
	return prob, nil
}

func (dist RedisWordProbabilityDist) ProbabilityOfTextAndLabel(text Text, label SpamLabel) (float64, error) {
	pTextGivenLabel, err := dist.ProbabilityOfTextGivenLabel(text, label)
	if err != nil {
		return 0.0, err
	}
	pLabel, err := dist.ProbabilityOfLabel(label)
	if err != nil {
		return 0.0, err
	}
	return pTextGivenLabel * pLabel, nil
}

func (dist RedisWordProbabilityDist) MostLikelyLabel(text Text) (SpamLabel, error) {
	spamProb, err := dist.ProbabilityOfTextAndLabel(text, SPAM)
	if err != nil {
		return false, err
	}
	hamProb, err := dist.ProbabilityOfTextAndLabel(text, HAM)
	if err != nil {
		return false, err
	}
	if spamProb > hamProb {
		return SPAM, nil
	} else {
		return HAM, nil
	}
}
